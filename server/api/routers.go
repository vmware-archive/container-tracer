package api

import (
	"github.com/gin-gonic/gin"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracerctx"
)

// map request path to logic
func NewRouter(t *ctx.Tracer) *gin.Engine {
	router := gin.Default()
	router.GET("/traces", t.SystemCallGet)
	router.GET("/traces/:id", t.SystemCallStatus)
	router.POST("/traces", t.SystemCallPost)
	router.DELETE("/traces/:id", t.SystemCallDelete)
	router.GET("/pods", t.LocalPodsGet)
	router.GET("/trace-hooks", t.TraceHooksGet)
	return router
}
