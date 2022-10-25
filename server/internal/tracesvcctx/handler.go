// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Frontend handlers of the trace-kube REST API.
 */
package tracekubectx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all pods, running on the local node
func (t *TraceKube) PodsGet(c *gin.Context) {
	c.JSON(http.StatusOK, "{}")
}

// get all trace hooks
func (t *TraceKube) TraceHooksGet(c *gin.Context) {
	c.JSON(http.StatusOK, "{}")
}

// get all trace sessions
func (t *TraceKube) TraceSessionGet(c *gin.Context) {
	//	id := c.Param("id")

	c.JSON(http.StatusOK, "{}")
}

// modify a trace session
func (t *TraceKube) TraceSessionPut(c *gin.Context) {
	//	id := c.Param("id")

	c.JSON(http.StatusOK, "{}")
}

// create a trace session
func (t *TraceKube) TraceSessionPost(c *gin.Context) {
	c.JSON(http.StatusOK, "{}")
}

// delete a trace session
// if id == "all", all trace sessions are deleted and trace subsystems are reseted
func (t *TraceKube) TraceSessionDel(c *gin.Context) {

	c.JSON(http.StatusOK, "{}")
}
