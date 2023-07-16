package configs

import (
	"github.com/gin-gonic/gin"
	middleware "github.com/nathanieiav/project-skripsi/routes/middlewares"
	"github.com/spf13/viper"
)

type routerGroup struct {
	router *gin.Engine
}

func RouterGroup() *gin.Engine {
	allowedHosts := viper.GetString("ALLOWED_HOSTS")

	router := routerGroup{
		router: gin.Default(),
	}

	router.router.SetTrustedProxies([]string{allowedHosts})
	router.router.Use(middleware.CORSMiddleware())

	v1 := router.router.Group("/v1/api")
	router.RegisteredRoute(v1)

	return router.router
}
