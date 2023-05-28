package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samridhi/golang-jwt-auth/routes"
)
func main(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT") //get env variable
	if port==""{
		port="8000"
	}
	router := gin.New() //creates a new router
	router.Use(gin.Logger()) //adds the built-in gin.Logger() middleware to a Gin router, used to perform common tasks like logging, authentication, and error handling
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	//gin.H(map[string]interface{})  make your code more concise and readable, especially when constructing complex JSON responses
	router.GET("/api-1",func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for api-1"})
	})
	router.GET("/api-2", func(c *gin.Context){
       c.JSON(200, gin.H{"success":"Access granted for api-2"})
	})
	router.Run(":" + port)

}