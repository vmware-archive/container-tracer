package api

import (
	"github.com/gin-gonic/gin"
)

// map request path to logic
func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/traces", SystemCallGet)
	router.POST("/traces", SystemCallPost)

	return router
}
