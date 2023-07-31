package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nathanieiav/project-skbackend/controllers"
)

func UserRoute(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.GET("/", controllers.UserController)
}
