package router

import (
	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func GetRouter() *gin.Engine {
	return r
}