package main

import (
	"os"

	routes "github.com/ShabnamHaque/go-jwt/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		//make sure correct port is present
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router) //points to the AuthRoutes func in router package
	routes.UserRoutes(router)

	//create the API(s)
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for API-1"})
	})

	//create the second API.
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for API-2"})
	})
	router.Run(":" + port)

}
