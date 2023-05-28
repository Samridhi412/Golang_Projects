package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samridhi/golang-jwt-auth/controllers"
	"github.com/samridhi/golang-jwt-auth/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate()) //put auth on users route as cant access without logging in
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}