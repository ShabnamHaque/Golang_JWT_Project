package routes

import (
	controller "github.com/ShabnamHaque/go-jwt/controllers"
	"github.com/gin-gonic/gin"
)

//func to authenticate user data
//post because it posts data fed into the database (client->server)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
