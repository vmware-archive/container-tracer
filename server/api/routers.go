package api

import (
	"github.com/gin-gonic/gin"
)

// map request path to logic
func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/traces", SystemCallGet)
	router.GET("/traces/:id", SystemCallStatus)
	router.POST("/traces", SystemCallPost)
	router.DELETE("/traces/:id", SystemCallDelete)

	return router
}
