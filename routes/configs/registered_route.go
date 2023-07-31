package configs

import (
	"github.com/gin-gonic/gin"
	"github.com/nathanieiav/project-skbackend/routes"
)

func (routerGroup routerGroup) RegisteredRoute(rg *gin.RouterGroup) {
	routes.UserRoute(rg)
}
