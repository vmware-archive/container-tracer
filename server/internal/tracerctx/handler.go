// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Enyinna Ochulor <eochulor@vmware.com>
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Backend handlers of the tracer REST API.
 */
package tracerctx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all pods, running on the local node
func (t *Tracer) LocalPodsGet(c *gin.Context) {
	if e := t.pods.Scan(); e != nil {
		c.JSON(http.StatusInternalServerError, e.Error())
	}
	cdb := t.pods.Get()
	if cdb != nil && len(*cdb) > 0 {
		c.JSON(http.StatusOK, cdb)
	} else {
		c.JSON(http.StatusOK, "{}")
	}
}

// get all trace hooks
func (t *Tracer) TraceHooksGet(c *gin.Context) {
	h := t.hooks.Get()
	if h != nil && len(*h) > 0 {
		c.JSON(http.StatusOK, h)
	} else {
		c.JSON(http.StatusOK, "{}")
	}
}

// get all trace sessions
func (t *Tracer) TraceSessionGet(c *gin.Context) {
	id := c.Param("id")

	if resp, err := t.getSession(&id, false); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		if resp != nil && len(*resp) > 0 {
			c.JSON(http.StatusOK, *resp)
		} else {
			c.JSON(http.StatusOK, "{}")
		}
	}
}

// modify a trace session
func (t *Tracer) TraceSessionPut(c *gin.Context) {
	var s sessionChange
	id := c.Param("id")

	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := t.changeSession(&id, &s); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, "{}")
}

// create a trace session
func (t *Tracer) TraceSessionPost(c *gin.Context) {
	var s sessionNew
	resp := traceSessionInfo{}

	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if id, err := t.newSession(&s); err != nil {
		c.JSON(http.StatusNotFound, err)
	} else if info, err := t.getSessionInfo(id); err == nil {
		resp = *info
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

// delete a trace session
func (t *Tracer) TraceSessionDel(c *gin.Context) {
	id := c.Param("id")

	if err := t.destroySession(&id); err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, "{}")
}
