// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * trace-kube service REST API.
 */
package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracesvcctx"
)

var (
	apiVersion = "v1"
)

// map request path to logic
func NewRouter(t *ctx.TraceKube) *gin.Engine {
	router := api.Router.SetupRouter()
	router.GET("/"+apiVersion+"/pods", t.ProxyAllMap)
	router.GET("/"+apiVersion+"/trace-hooks", t.ProxyAnyMap)
	router.POST("/"+apiVersion+"/trace-session", t.ProxyAllMap)
	router.GET("/"+apiVersion+"/trace-session/:id", t.ProxyAllMap)
	router.PUT("/"+apiVersion+"/trace-session/:id", t.ProxyAllMap)
	router.DELETE("/"+apiVersion+"/trace-session/:id", t.ProxyAllMap)
	return router
}
