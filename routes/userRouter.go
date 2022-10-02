package routes

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/MahadevMaahi/gin-jwt-authentication/middlewares"
	controller "github.com/MahadevMaahi/gin-jwt-authentication/controllers"
)

func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Autheticate())
	routes.GET("/users", controller.GetUsers())
	routes.GET("/users/:user_id", controller.GetUser())
}