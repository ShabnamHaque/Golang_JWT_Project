package routes

import (
	controller "github.com/ShabnamHaque/go-jwt/controllers"
	"github.com/ShabnamHaque/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	//use middleware because the two func below need to be protected,accessed only via a token.
	
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
