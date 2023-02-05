Directory with trace-cruncher scripts, used to run ftrace on containers. All user callable scripts must:
- Have prefix **trace_**, in order to be auto discovered.
- Support at least these arguments:
    - **--pid** : list of Process IDs to be traced, a mandatory argument.
    - **--instance** : Name of the trace instance used for tracing, optional argument.
    - **--time** : Duration of the trace in milliseconds, optional argument.
    - **--describe** : Return a user description of the script.
- The scripts run in blocking mode and must support graceful termination with the signals **SIGUSR1** or **SIGINT**.

This common functionality is implemented in `tc_base.py`, it can be reused by scripts by inheriting `class tracer`.