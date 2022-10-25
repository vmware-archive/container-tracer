// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * trace-kube service REST API.
 */
package api

import (
	"github.com/gin-gonic/gin"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracesvcctx"
)

var (
	apiVersion = "v1"
)

// map request path to logic
func NewRouter(t *ctx.TraceKube) *gin.Engine {
	router := gin.Default()
	router.GET("/"+apiVersion+"/pods", t.PodsGet)
	router.GET("/"+apiVersion+"/trace-hooks", t.TraceHooksGet)
	router.POST("/"+apiVersion+"/trace-session", t.TraceSessionPost)
	router.GET("/"+apiVersion+"/trace-session/:id", t.TraceSessionGet)
	router.PUT("/"+apiVersion+"/trace-session/:id", t.TraceSessionPut)
	router.DELETE("/"+apiVersion+"/trace-session/:id", t.TraceSessionDel)
	return router
}
