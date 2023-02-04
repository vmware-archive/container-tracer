# tracer-node overview
![tracer-node](tracer-node.png)
## Architecture
The `tracer-node` module provides the main container-tracer functionality for container tracing.
It has the following main components:  
- A REST API, used to interact with the tracer on that node. Look at [REST API](container-tracer-api.md)
  for the API description.  
- Logic for auto-discovery of all pods running on the node. Two different approaches are used
  for this auto-discovery:  
    - Using the CRI API. This is the preferred approach, when container-tracer runs in a Kubernetes context.  
    - Using the information from the `/proc` file system on the host. If the CRI API is
    not available, this logic is used.  
- An in-memory database with all pods and containers running on the node. For each container,
  a list of PIDs is stored into the database, as seen in the host PID namespace.  
- A list of [trace-hooks](container-tracer-hooks.md), available in the `tracer-node`.  
- An in-memory database with configured trace sessions. A trace session is a set of containers,
  trace hook and trace parameters that has a state - running or stopped. When running, the trace
  hook is attached to the specified containers.  
- Open Telemetry trace exporters, used to export the output of running trace sessions to an
  external database.
## Parameters
On startup, `tracer-node` checks for specific environment variables and accepts these input arguments:  
- `--address` or `TRACER_API_ADDRESS`:  IP address and port in format IP:port, used to listen
   for incoming REST API requests. By default `:8080` is used.  
- `--node-name` or `TRACER_NODE_NAME`: The name of the node, that runs this `tracer-node` instance. Used
for logging purposes and is included in the traces collected on this node. There is no default value.  
- `--pod-name` or `TRACER_POD_NAME`: The name of this `tracer-node` pod, used to verify the correct
CRI endpoint.  
- `--cri-endpoint` or `TRACER_CRI_ENDPOINT`: A specific CRI endpoint that must be used for CRI API.
There is no default value, if it is not set - an auto discovery logic is used.  
- `--run-path` or `TRACER_RUN_PATHS`: The path to the run directories of the host, to search for cri
endpoints. By default `/run` and `/var/run` are used, but usually when running in a container, the host
run paths are mounted on custom locations. These are used to auto discover the endpoint of the CRI API,
if no specific CRI endpoint is specified. 
- `--procfs-path` or `TRACER_PROCFS_PATH`: The path to the host `/proc` file system mount point.
By default it is `/proc`, but usually when running in a container the host `/proc` is mounted on
a custom location. 
- `--sysfs-path` or `TRACER_SYSFS_PATH`: The path to the host `/sys` file system mount point.
By default it is `/sys`, but usually when running in a container, the host `/sys` is mounted on
a custom location.  
- `--use-procfs` or `TRACER_FORCE_PROCFS`: Force the use of `/proc` of the host for auto-discovery
of the containers, running on the local node. Not set by default. The default logic is using
the CRI API if it is available. If the CRI API is not accessible, fail back to the logic that gets
this information from the `/proc` file system.  
- `--jaeger-endpoint` or `TRACER_JEAGER_ENDPOINT`: The URL of the jaeger endpoint service, used to send
the collected traces. Can be set to `auto`, which triggers the default logic - search for
`jaeger-collector` service that exposes port `14268` and use `http://jaeger-collector:14268/api/traces`
as an endpoint to jaeger.  
- `--verbose` or `TRACE_KUBE_VERBOSE`: Dump more detailed logs, disabled by default.  

If both input argument and environment variable for a same setting exist, only the input argument is taken.