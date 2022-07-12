![Mind map diagram](https://gitlab.eng.vmware.com/opensource/tracecruncher-api/-/raw/main/Mind_Map.png)

Components to work on for the traceCruncher API:


a. To start, the focus is on having an API endpoint for system calls.  One of the examples that currently exists in trace-cruncher will be used.

b. Parameter/s: ID of container or part of the container name, custom string to pass to OpenTelemetry and this should identify the session.

c. Find the containers that are running that have the name specified by the user and get the PIDs of the containers that are running in the cluster.

d. Maintain a database that stores the tracing ID that is generated when a tracing resource is created. This ID is used in starting, stopping and destroying of the tracing resource.
    
    - Database: store parameters related to the creation of a tracing resource; tracing ID, entire information received with the create tracing resource.

e. Implement a script in python that reads the ftrace buffer because each tracing session will be a different ftrace buffer. Data will be sent to OpenTelemetry.

f. Verbs that will be needed: create,start, stop, status, destroy.
    
    - POST(Create): will enter a tracing resource in the API database.
    
    - PUT(Start): will run given traceCruncher logic with given parameters(payload). Also, start will run additional script from “e”.
    
    - PUT(Stop): destroys ftrace tracing session. Ftrace tracing can be restarted with start. (decide if the script will run in the background?). Kill ftrace buffer reading script.
    
    - DELETE(Destroy): stop the ftrace session if it is not already stopped. Remove the tracing resource from the API database.
    
    - GET(Status): look in the internal database and print the status metadata that is stored for the tracing resource.

## Docker

tracecruncher-api is installable via Docker. The installation includes the trace-cruncher library.

To install run the following from the `server` directory:
    docker build . -t tracecruncher-api-image -f Dockerfile
It's recommended to run the installer with the `--squash` flag in order to reduce the size of the final image.

To run the container, use the `--priviliged` flag to give trace-cruncher kernel access and ensure you are publishing the container port:
    docker run --privileged -p 8080:8080 --name tracecruncher-api -it tracecruncher-api-image
