package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nathanieiav/project-skripsi/controllers"
)

func UserRoute(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.GET("/", controllers.UserController)
}
