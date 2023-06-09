package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samridhi/golang-jwt-auth/helpers"
)
func Authenticate() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == ""{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No authorization header provided")})
			ctx.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != ""{
             ctx.JSON(http.StatusInternalServerError, gin.H{"error":err})
			 ctx.Abort()
			 return
		}
		ctx.Set("email",claims.Email)
		ctx.Set("first_name",claims.First_name)
		ctx.Set("last_name",claims.Last_name)
		ctx.Set("uid",claims.Uid)
		ctx.Set("user_type",claims.User_type)
		ctx.Next()

	}
}