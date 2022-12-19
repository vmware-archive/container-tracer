#!/usr/bin/env python3
"""
SPDX-License-Identifier: GPL-2.0-or-later
Copyright 2022 VMware Inc, Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

Manager of trace helper programs located in the same directory:
 - auto discovers available trace programs
 - gets their description and arguments
 - run a trace program
"""

import os, subprocess, signal, sys
import argparse, string, random
from pathlib import Path
import tracecruncher.ftracepy as ft

scripts_dir="./"
scripts_prefix = "trace_"
instance_prefix = "kube_"
instance_name_len = 16 - len(instance_prefix)
instance_name_chars = string.ascii_letters + string.digits
max_ftrace_retries = 10
envProc = "TRACER_PROCFS_PATH"
envSys = "TRACER_SYSFS_PATH"
procDefPath = "/proc"
sysDefPath = "/sys"

description="Manager of trace scripts in the current directory"

def get_scripts():
    for file in os.listdir(scripts_dir):
      if file.startswith(scripts_prefix):
        print(Path(file).stem)
    exit(0)

def run_script(name, arguments):
    signal.signal(signal.SIGUSR1, signal.SIG_IGN)
    for file in os.listdir(scripts_dir):
      if file.startswith(name + "."):
        try:
            output = subprocess.run(["./" + file] + arguments, capture_output=True, universal_newlines = True, start_new_session=True)
            print(output.stdout, flush=True)
            print(output.stderr, file=sys.stderr, flush=True)
        except KeyboardInterrupt:
            pass
    exit(0)

def run_trace(name, arguments):
    instance = None
    retries = max_ftrace_retries
    while not instance and retries > 0:
      try:
        rname =  ''.join(random.choice(instance_name_chars) for _ in range(instance_name_len))
        iname = instance_prefix + rname
        instance = ft.create_instance(tracing_on=False, name=iname)
      except:
        retries-=1
        pass
    if not instance:
      raise RuntimeError("Failed to create a trace instance")
    print(ft.dir()+"/instances/"+iname+"/trace_pipe", flush=True)
    run_script(name, arguments + ["--instance", iname])

def reset_ftrace():
    instances = [f.path for f in os.scandir(ft.dir()+"/instances") if f.is_dir()]
    for i in list(instances):
      if (Path(i).stem.startswith(instance_prefix)):
        os.rmdir(i)
    exit(0)

def set_ftrace_dir():
    eproc = os.environ.get(envProc)
    esys = os.environ.get(envSys)
    # Return if no custom /proc or /sys mount points are passed
    if not eproc and not esys:
        return
    if not eproc:
        eproc = procDefPath
    mfile = open(eproc + '/mounts', 'r')
    lines = mfile.readlines()
    ftraceMount=""
    debugMount=""
    # Look for tracefs or debugfs in /proc/mounts
    for l in lines:
        w = l.split(" ")
        if len(w) < 3:
            continue
        if w[2] == 'debugfs'and os.path.isdir(w[1]):
            debugMount = w[1]
        if w[2] == 'tracefs'and os.path.isdir(w[1]):
            ftraceMount = w[1]
            break
    # If tracefs is not found, check in debugfs
    if ftraceMount == "" and debugMount != "":
        fm = debugMount + 'tracing'
        if os.path.isdir(fm):
            ftraceMount = fm
    if ftraceMount != "":
        if esys and ftraceMount.startswith('/sys'):
            fm = ftraceMount.removeprefix('/sys')
            ftraceMount = esys + fm
        ft.set_dir(ftraceMount)
    else:
        print("Failed to find ftrace mount point", file=sys.stderr, flush=True)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=description)
    parser.add_argument('-g', '--get-all', action='store_true', dest='get_all',
                        help="Get available scripts")
    parser.add_argument('-d', '--describe', nargs=1, dest='get_desc',
                        help="Get description of a script")
    parser.add_argument('-c', '--clear', action='store_true', dest='reset',
                        help="Reset ftrace subsystem to default")
    parser.add_argument('-r', '--run', dest='script', nargs=1, help="Name of a trace script to run")
    parser.add_argument('-a', '--args', dest='arguments',  nargs=1,
                        help="Arguments of a trace script")

    set_ftrace_dir()

    args = parser.parse_args()

    if args.get_all:
      get_scripts()
    if args.get_desc:
      run_script(args.get_desc[0], ["--describe"])
    if args.reset:
        reset_ftrace()
    if args.script:
        run_trace(args.script[0], list(args.arguments[0].split(" ")))
