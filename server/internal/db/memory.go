package db

// trace represents data about a tracing resource
// ID of container or part of the container name, custom string
// to pass to OpenTelemetry and this should identify the session
type Trace struct {
	ID        string `json:"id"`
	Container string `json:"name"`
}

// Sample tracing resources to start with
var Traces = []Trace{
	{ID: "1", Container: "Alpha"},
	{ID: "2", Container: "Beta"},
	{ID: "3", Container: "Gamma"},
}

// sample POST

/*
curl http://localhost:8080/traces \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","name": "Delta"}'
*/

// sample GET

/*
curl http://localhost:8080/traces \
    --header "Content-Type: application/json" \
    --request "GET"
*/
