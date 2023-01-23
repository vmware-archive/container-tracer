# trace hooks
`Container-tracer` uses a set of trace hooks to interact with the tracing infrastructure of
the Linux kernel. The trace hooks are responsible to configure a tracing session with
given context in the kernel and to return a location of a text file, where traces of
this session are recorded. Usually, this tracing context is a set of PIDs that will be
traced and input arguments, specific to the trace-hook. The trace hooks are located in
`trace-hooks` directory and are organized into sub-directories. Each sub-directories is
responsible for a specific trace sub-system of the Linux kernel. There is at least one
executable file in each of these sub-directories, called `manager`, which manages all
trace hooks located in that sub-directory. The `manager` must accept at least these input
arguments:  
 - **--get-all** : Return list of all user callable trace hooks.
 - **--describe <trace hook name>** : Get user description of given trace hook.
 - **--clear** : Reset to default the trace sub-system of the Linux kernel.
 - **--run <trace hook name>** : Run a trace hook, in blocking mode. In case of an error, an error message must
   be printed on the standard error output and the hook must return, without starting any trace session.
   If everything is OK, only the full path to a file where traces are collected must be printed on the
   standard output, no prints on the standard error. The trace hook blocks this instance of the `manager`
   during the trace session. The trace session stops when this instance of `manager` receives a **SIGINT**
   signal.
 - **--args <trace hook arguments>** : Arguments that will be passed to the trace hook.

These environment variables can be used to set system specific configuration to trace hooks, the `manager`
must read them and apply this configuration:  
- **TRACER_PROCFS_PATH**: Mount location of the host **/proc** file system.
   If not set, the default **/proc** is used.
- **TRACER_SYSFS_PATH**: Mount location of the host **/sys** file system.
   If not set, the default **/sys** is used.

`Container-tracer` uses `manager` to auto-discover and run available trace hooks. New types of
trace hooks, to a different tracing subsystem, can be added easily by creating a new sub-directory
in `trace-hooks` and implementing `manager` for them.