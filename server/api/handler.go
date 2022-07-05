package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/db"
)

// create a new tracing reource
func SystemCallPost(c *gin.Context) {
	var newTrace db.Trace

	// bind the received JSON to newTrace
	if err := c.BindJSON(&newTrace); err != nil {
		return
	}

	// add new trace to the Traces slice
	db.Traces = append(db.Traces, newTrace)
	c.IndentedJSON(http.StatusCreated, newTrace)
}

// get all the the current traces and respond as JSON
func SystemCallGet(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, db.Traces)
}
