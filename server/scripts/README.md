These subdirectories contain the implementation of various trace scripts, related to a specific tracing framework.

In order to manage the scripts in a unified way, each subdirectory must have a `scripts_manager.py` which supports at least these options:
 - **-g, --get-scripts** : Return list of all autodiscovered user callable scripts.
 - **-d, --describe-script** : Get user description of given script.
 - **-c, --clear** : Reset to default the ftrace subsystem.
 - **-r, --run** : Run a trace script, in blocking mode. A full path to a file where traces are collected is printed on the standard output.
 - **-a, --args** : Arguments that will be passed to the script.
