package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/MahadevMaahi/gin-jwt-authentication/controllers"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("users/signup", controller.Signup())
	routes.POST("users/login", controller.Login())
}