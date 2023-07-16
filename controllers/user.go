package controllers

import "github.com/gin-gonic/gin"

func UserController(context *gin.Context) {
	context.String(200, "hellooooo")
}
