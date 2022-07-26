#!/usr/bin/env python3

"""
SPDX-License-Identifier: CC-BY-4.0

Copyright 2022 VMware Inc, Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
"""

import os, subprocess, signal
import argparse, string, random
from pathlib import Path
import tracecruncher.ftracepy as ft

scripts_dir="./"
scripts_prefix = "trace_"
instance_prefix = "kube_"
instance_name_len = 16 - len(instance_prefix)
instance_name_chars = string.ascii_letters + string.digits
max_ftrace_retries = 10

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
            print(output.stdout)
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
    print(ft.dir()+"/instances/"+iname+"/trace")
    run_script(name, arguments + ["--instance", iname])

def reset_ftrace():
    instances = [f.path for f in os.scandir(ft.dir()+"/instances") if f.is_dir()]
    for i in list(instances):
      if (Path(i).stem.startswith(instance_prefix)):
        os.rmdir(i)
    exit(0)

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

    args = parser.parse_args()

    if args.get_all:
      get_scripts()
    if args.get_desc:
      run_script(args.get_desc[0], ["--describe"])
    if args.reset:
        reset_ftrace()
    if args.script:
        run_trace(args.script[0], list(args.arguments[0].split(" ")))
