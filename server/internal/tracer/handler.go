package tracer

import (
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
