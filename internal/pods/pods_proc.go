// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Discover containers running on the local node, using information from the /proc file system.
 * This logic was originally implemented in python by Yordan Karadzhov <y.karadz@gmail.com>
 *
 * Using /proc file system has limitations in Kubernetes context. I couldn't find reliable way to get
 * Pod -> Containers relation, so the logic considers that all tasks inside a Pod are part of a single
 * container with name "unknown".
 */
package pods

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	defaultContainer = "unknown"
)

type podProc struct {
	path         string
	podb         map[string]*pod
	top_utsns_fd int
	top_utsns_id int
	ppid         int
	pids         []int
}

func getProcDiscover(procfsPath *string) (podsDiscover, error) {
	ctr := podProc{
		path: *procfsPath,
		pids: make([]int, 0),
		podb: make(map[string]*pod),
	}

	if id, err := ctr.getNSinum(1, "uts"); err != nil {
		return nil, err
	} else {
		ctr.top_utsns_id = id
	}

	if f, err := os.Open(fmt.Sprintf("%s/1/ns/uts", ctr.path)); err != nil {
		return nil, err
	} else {
		ctr.top_utsns_fd = int(f.Fd())
	}

	print("\nUsing PROC for pods discovery at ", ctr.path, "\n")
	return ctr, nil
}

func (p *podProc) getNSinum(pid int, ns string) (int, error) {
	if name, err := os.Readlink(fmt.Sprintf("%s/%d/ns/%s", p.path, pid, ns)); err == nil {
		f := func(c rune) bool {
			return c == '[' || c == ']'
		}
		fields := strings.FieldsFunc(name, f)
		if len(fields) != 2 || fields[0] != fmt.Sprintf("%s:", ns) {
			return -1, fmt.Errorf("Broken name space \"%s\" id: \"%s\"", ns, name)
		}
		return strconv.Atoi(fields[1])
	} else {
		return -1, err
	}
}

func (p *podProc) nsGetPodName(pid int) (*string, error) {
	if f, err := os.Open(fmt.Sprintf("%s/%d/ns/uts", p.path, pid)); err == nil {
		defer f.Close()
		if err := unix.Setns(int(f.Fd()), unix.CLONE_NEWUTS); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	utsname := &unix.Utsname{}
	if err := unix.Uname(utsname); err != nil {
		return nil, err
	}

	if err := unix.Setns(p.top_utsns_fd, unix.CLONE_NEWUTS); err != nil {
		return nil, err
	}

	name := string(utsname.Nodename[:bytes.IndexByte(utsname.Nodename[:], 0)])
	return &name, nil
}

func (p *podProc) compareNS(pid1, pid2 int, ns string) (bool, error) {
	var ns1, ns2 int
	var err error

	if ns1, err = p.getNSinum(pid1, ns); err != nil {
		return false, err
	}
	if ns1, err = p.getNSinum(pid1, ns); err != nil {
		return false, err
	}

	return ns1 == ns2, nil
}

func (p *podProc) getPodInfo(ppid, pid int) error {
	// Check if the given PID is in different pid name space from its parrent
	if b, err := p.compareNS(ppid, pid, "pid"); err != nil {
		return err
	} else if b {
		return nil
	}

	// Check if the PID is in different uts name space from process 1
	if uts, err := p.getNSinum(pid, "uts"); err != nil {
		return err
	} else if uts == p.top_utsns_id {
		return nil
	}

	// Get the name of the pod and all children of the given PID
	if name, err := p.nsGetPodName(pid); err == nil {
		if v, ok := p.podb[*name]; ok {
			t := v.Containers[defaultContainer].Tasks
			t = append(t, pid)
		} else {
			p.podb[*name] = &pod{
				Containers: map[string]*Container{
					defaultContainer: &Container{
						Id:    &defaultContainer,
						Pod:   name,
						Tasks: []int{pid},
					},
				},
			}
		}
		if err := p.getChildren(pid, true); err != nil {
			return err
		}
		if len(p.pids) > 0 {
			t := p.podb[*name].Containers[defaultContainer].Tasks
			t = append(t, p.pids...)
		}
	} else {
		return err
	}

	return nil
}

func (p *podProc) getChildrenPids(pid int) error {
	if file, err := os.Open(fmt.Sprintf("%s/%d/task/%d/children", p.path, p.ppid, pid)); err == nil {
		defer file.Close()
		scan := bufio.NewScanner(file)
		scan.Split(bufio.ScanWords)
		for scan.Scan() {
			if i, e := strconv.Atoi(scan.Text()); e == nil {
				p.pids = append(p.pids, i)
			}
		}
	} else {
		return err
	}
	return nil
}

func (p *podProc) procWalk(path string, di fs.DirEntry, err error) error {
	if pid, err := strconv.Atoi(di.Name()); err == nil {
		if di.IsDir() || p.ppid == pid {
			return nil
		}
		p.getChildrenPids(pid)
	}

	return nil
}

func (p *podProc) getChildren(pid int, threads bool) error {
	var err error
	p.ppid = pid
	p.pids = p.pids[:0]
	if err = p.getChildrenPids(pid); err != nil {
		return err
	}

	if !threads {
		return nil
	}
	if err := filepath.WalkDir(fmt.Sprintf("%s/%d/task/", p.path, p.ppid), p.procWalk); err != nil {
		return err
	}

	return nil
}

func (p *podProc) walkChildren(ppid, pid int) error {

	if err := p.getChildren(pid, false); err != nil {
		return err
	}
	walk := make([]int, len(p.pids))
	copy(walk, p.pids)

	if ppid > 1 {
		if err := p.getPodInfo(ppid, pid); err != nil {
			return err
		}
	}
	for _, pr := range walk {
		if err := p.walkChildren(pid, pr); err != nil {
			return err
		}
	}

	return nil
}

func (p podProc) podScan() (*map[string]*pod, error) {
	// Reset the pods databse
	p.podb = make(map[string]*pod)

	// Dsicover the children of process 1
	if err := p.walkChildren(0, 1); err != nil {
		return nil, err
	}

	return &p.podb, nil
}
