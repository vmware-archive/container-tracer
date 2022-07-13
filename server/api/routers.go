package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracer"
)

// map request path to logic
func NewRouter(t *tracer.Tracer) *gin.Engine {
	router := gin.Default()
	router.GET("/traces", t.SystemCallGet)
	router.GET("/traces/:id", t.SystemCallStatus)
	router.POST("/traces", t.SystemCallPost)
	router.DELETE("/traces/:id", t.SystemCallDelete)
	router.GET("/containers", t.LocalContainersGet)

	return router
}
