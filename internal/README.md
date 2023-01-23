# container-tracer implementation
Internal implementation of [tracer-node](../docs/tracer-node.md) and [tracer-svc](../docs/tracer-svc.md) modules.

## tracer-node internals
- **tracerctx**: Main logic of `tracer-node`. Initialize all other sub-modules and maintain the runtime
  context of this `tracer-node` instance. Implementation of the REST API handlers. Database and logic for
  running trace sessions.
- **tracehook**: Logic for working with trace hooks - auto discovery available hooks; run and terminate
   a hook as part of a trace session, read standard output and error of a trace hook instance.
- **pods**: Database and logic for auto-discovery of PODs and containers, running on the local system.
- **logger**: Implementation of trace exporters to external databases, using Open Telemetry SDK.

## tracer-svc internals
- **tracesvcctx**: Main logic of `tracer-svc`. Maintain the runtime context of this `tracer-svc`
  instance. Database and logic for auto-discovery of `tracer-node` instances, running on the cluster.
  Implementation of the REST API handlers and broadcast proxy.
