![Mind map diagram](https://gitlab.eng.vmware.com/opensource/tracecruncher-api/-/blob/main/Mind_Map.png)*Mind map*

Components to work on for the traceCruncher API:


a. To start, the focus is on having an API endpoint for system calls.  One of the examples that currently exists in trace-cruncher will be used.

b. Parameter/s: ID of container or part of the container name, custom string to pass to OpenTelemetry and this should identify the session.

c. Find the containers that are running that have the name specified by the user and get the PIDs of the containers that are running in the cluster.

d. Maintain a database that stores the tracing ID that is generated when a tracing resource is created. This ID is used in starting, stopping and destroying of the tracing resource.
    - Database: store parameters related to the creation of a tracing resource; tracing ID, entire information received with the create tracing resource.

e. Implement a script in python that reads the ftrace buffer because each tracing session will be a different ftrace buffer. Data will be sent to OpenTelemetry.

f. Verbs that will be needed: create,start, stop, status, destroy.
    - Create: will create entry in the API database.
    - Start: will run given traceCruncher script with given parameters. Also, start will run additional script from “e”.
    - Stop: destroys ftrace tracing session. Ftrace tracing can be restarted with start. (decide if the script will run in the background?). Kill ftrace buffer reading script.
    - Destroy: stop the ftrace session if it is not already stopped. Remove the tracing resource from the API database.
    - Status: Look in the internal database and print the status metadata that is stored for the tracing resource.
