// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Enyinna Ochulor <eochulor@vmware.com>
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Backend handlers of the tracer REST API.
 */
package tracerctx

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/db"
)

// create a new tracing reource
func (t *Tracer) SystemCallPost(c *gin.Context) {
	var newTrace db.Trace

	// bind the received JSON to newTrace
	if err := c.BindJSON(&newTrace); err != nil {
		return
	}

	// add new trace to the Traces slice
	db.Traces = append(db.Traces, newTrace)
	c.IndentedJSON(http.StatusCreated, newTrace)
}

// get all the current traces and respond as JSON
func (t *Tracer) SystemCallGet(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.Traces)
}

// get the status of a specific trace by ID
func (t *Tracer) SystemCallStatus(c *gin.Context) {
	id := c.Param("id")

	// loop through the Traces slice to find a match
	for _, a := range db.Traces {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "trace not found"})
}

// remove/stop an existing trace using a Trace ID
func (t *Tracer) SystemCallDelete(c *gin.Context) {
	id := c.Param("id")

	// loop through the Traces slice to find and remove the
	// Trace with the specified ID
	for i, a := range db.Traces {
		if a.ID == id {
			out := a
			db.Traces = removeTrace(db.Traces, i)
			c.IndentedJSON(http.StatusOK, out)
		}
	}
}

// this is a utility function to delete an element from the Traces slice
func removeTrace(trace_arr []db.Trace, i int) []db.Trace {
	trace_arr[i] = trace_arr[len(trace_arr)-1]
	return trace_arr[:len(trace_arr)-1]
}

// get all pods, running on the local node
func (t *Tracer) LocalPodsGet(c *gin.Context) {
	if e := t.pods.Scan(); e != nil {
		c.IndentedJSON(http.StatusInternalServerError, e)
	}
	cdb := t.pods.Get()
	if cdb != nil && len(*cdb) > 0 {
		if j, e := json.Marshal(cdb); e != nil {
			c.IndentedJSON(http.StatusInternalServerError, e)
		} else {
			c.IndentedJSON(http.StatusOK, string(j))
		}
	} else {
		c.IndentedJSON(http.StatusOK, "{}")
	}
}

// get all trace hooks
func (t *Tracer) TraceHooksGet(c *gin.Context) {
	h := t.hooks.Get()
	if h != nil && len(*h) > 0 {
		if j, e := json.Marshal(h); e != nil {
			c.IndentedJSON(http.StatusInternalServerError, e)
		} else {
			c.IndentedJSON(http.StatusOK, string(j))
		}
	} else {
		c.IndentedJSON(http.StatusOK, "{}")
	}
}
