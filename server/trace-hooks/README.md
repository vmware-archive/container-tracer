These subdirectories contain implementation of various trace helpers, related to a specific tracing framework.

In order to manage the helpers in a unified way, each subdirectory must have an application called `manager` which supports at least these options:
 - **--get-all** : Return list of all user callable helpers.
 - **--describe** : Get user description of given helper.
 - **--clear** : Reset to default the trace configuration.
 - **--run** : Run a trace helper, in blocking mode. A full path to a file where traces are collected should be printed on the standard output.
 - **--args** : Arguments that will be passed to the helper.
