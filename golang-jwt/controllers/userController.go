package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/samridhi/golang-jwt-auth/database"
	"github.com/samridhi/golang-jwt-auth/helpers"
	"github.com/samridhi/golang-jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func VerifyPassword(userPassword string, providedPassword string) (bool,string){
         err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
		 check := true
		 msg := ""
		 if err != nil{
			msg = fmt.Sprintf("email or password is incorrect")
			check = false
		 }
		 return check, msg
}
func HashPassword(password string) string{
   bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
   if err != nil{
     log.Panic(err)
   }
   return string(bytes)
}

func Signup() gin.HandlerFunc{
    return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := ctx.BindJSON(&user); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(c, bson.M{"email":user.Email})
		defer cancel()
		if err != nil{
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"error ocurred while checking for user"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password
        count, err = userCollection.CountDocuments(c, bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil{
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"erorr":"error ocurred while checking phone number"})
		}
		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone already exists in records"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		resultInsertionNumber, insertErr := userCollection.InsertOne(c, user)
		if insertErr != nil{
			msg := fmt.Sprintf("User Item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc{
   return func(ctx *gin.Context) {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var foundUser models.User
	if err := ctx.BindJSON(&user); err != nil{
        ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	err := userCollection.FindOne(c, bson.M{"email":user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"email or password is incorrect"})
		return
	}
	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		return
	}
	if foundUser.Email == nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"User not found"})
	}
	token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
   helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
   err = userCollection.FindOne(c, bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)
   if err != nil{
	ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
	return
   }
   ctx.JSON(http.StatusOK, foundUser)
   }
}



//admin can have access to users data
func GetUser() gin.HandlerFunc{
    return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")
		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		} //to check admin
		var c, cancel =context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		//mongo saves data as json
		err := userCollection.FindOne(c, bson.M{"user_id":userId}).Decode(&user)
        defer cancel()
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func GetUsers() gin.HandlerFunc{
   return func(ctx *gin.Context) {
	if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
	if err != nil || recordPerPage<1{
         recordPerPage = 10
	}
	page, err1 := strconv.Atoi(ctx.Query("page"))
	if err1 != nil || page<1{
		page = 1
	}
	//pagination
	startIndex := (page-1)*recordPerPage //skip and limit
	startIndex, err = strconv.Atoi(ctx.Query("startIndex"))
	//aggregation in mongodb
	matchStage := bson.D{{"$match", bson.D{{}}}}
	//without push u dont get access to data just the count
	//group documents by id, start finding counts of unique users, and get the users
	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"_id", "null"}}},
		{"total_count", bson.D{{"$sum" ,1}}}, 
		{"data", bson.D{{"$push", "$$ROOT"}}},
		}}} //$sum - sum all records, group by id and count
	//which datapoints should go to user, more readable like graphql
	projectStage := bson.D{
          {"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data",startIndex, recordPerPage}}}},
		  }}}
	result, err := userCollection.Aggregate(c, mongo.Pipeline{
		matchStage, groupStage, projectStage })
	defer cancel()
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"error ocurred while listing user items"})
	}
	var allUsers []bson.M
	if err = result.All(c, &allUsers); err != nil{
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, allUsers[0])


   }
}