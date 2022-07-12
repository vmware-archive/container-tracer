// SPDX-License-Identifier: Apache-2.0
// Copyright (C) 2020 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>

/*
 Discover containers running on the local node, using information from the /proc file system
 This logic was originally implemented in python by Yordan Karadzhov <y.karadz@gmail.com>
*/

package condb

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

type containerProc struct {
	condb        map[string]*container
	top_utsns_fd int
	top_utsns_id int
	ppid         int
	pids         []int
}

func getProcDiscover() (containersDiscover, error) {
	ctr := containerProc{
		pids:  make([]int, 0),
		condb: make(map[string]*container),
	}

	if id, err := getNSinum(1, "uts"); err != nil {
		return nil, err
	} else {
		ctr.top_utsns_id = id
	}

	if f, err := os.Open("/proc/1/ns/uts"); err != nil {
		return nil, err
	} else {
		ctr.top_utsns_fd = int(f.Fd())
	}

	return ctr, nil
}

func getNSinum(pid int, ns string) (int, error) {
	if name, err := os.Readlink(fmt.Sprintf("/proc/%d/ns/%s", pid, ns)); err == nil {
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

func (c *containerProc) nsGetContName(pid int) (*string, error) {
	if f, err := os.Open(fmt.Sprintf("/proc/%d/ns/uts", pid)); err == nil {
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

	if err := unix.Setns(c.top_utsns_fd, unix.CLONE_NEWUTS); err != nil {
		return nil, err
	}

	name := string(utsname.Nodename[:bytes.IndexByte(utsname.Nodename[:], 0)])
	return &name, nil
}

func compareNS(pid1, pid2 int, ns string) (bool, error) {
	var ns1, ns2 int
	var err error

	if ns1, err = getNSinum(pid1, ns); err != nil {
		return false, err
	}
	if ns1, err = getNSinum(pid1, ns); err != nil {
		return false, err
	}

	return ns1 == ns2, nil
}

func (c *containerProc) getContainerInfo(ppid, pid int) error {
	// Check if the given PID is in different pid name space from its parrent
	if b, err := compareNS(ppid, pid, "pid"); err != nil {
		return err
	} else if b {
		return nil
	}

	// Check if the PID is in different uts name space from process 1
	if uts, err := getNSinum(pid, "uts"); err != nil {
		return err
	} else if uts == c.top_utsns_id {
		return nil
	}

	// Get the name of the container and all children of the given PID
	if name, err := c.nsGetContName(pid); err == nil {
		if v, ok := c.condb[*name]; ok {
			v.Pids = append(v.Pids, pid)
		} else {
			c.condb[*name] = &container{
				Pids: []int{pid},
			}
		}
		if err := c.getChildren(pid, true); err != nil {
			return err
		}
		if len(c.pids) > 0 {
			c.condb[*name].Pids = append(c.condb[*name].Pids, c.pids...)
		}
	} else {
		return err
	}

	return nil
}

func (c *containerProc) getChildrenPids(pid int) error {
	if file, err := os.Open(fmt.Sprintf("/proc/%d/task/%d/children", c.ppid, pid)); err == nil {
		defer file.Close()
		scan := bufio.NewScanner(file)
		scan.Split(bufio.ScanWords)
		for scan.Scan() {
			if i, e := strconv.Atoi(scan.Text()); e == nil {
				c.pids = append(c.pids, i)
			}
		}
	} else {
		return err
	}
	return nil
}

func (c *containerProc) procWalk(path string, di fs.DirEntry, err error) error {
	if pid, err := strconv.Atoi(di.Name()); err == nil {
		if di.IsDir() || c.ppid == pid {
			return nil
		}
		c.getChildrenPids(pid)
	}

	return nil
}

func (c *containerProc) getChildren(pid int, threads bool) error {
	var err error
	c.ppid = pid
	c.pids = c.pids[:0]
	if err = c.getChildrenPids(pid); err != nil {
		return err
	}

	if !threads {
		return nil
	}
	if err := filepath.WalkDir(fmt.Sprintf("/proc/%d/task/", c.ppid), c.procWalk); err != nil {
		return err
	}

	return nil
}

func (c containerProc) walkChildren(ppid, pid int) error {

	if err := c.getChildren(pid, false); err != nil {
		return err
	}
	walk := make([]int, len(c.pids))
	copy(walk, c.pids)

	if ppid > 1 {
		if err := c.getContainerInfo(ppid, pid); err != nil {
			return err
		}
	}
	for _, p := range walk {
		if err := c.walkChildren(pid, p); err != nil {
			return err
		}
	}

	return nil
}

func (c containerProc) contScan() (*map[string]*container, error) {
	// Reset the containers databse
	c.condb = make(map[string]*container)
	// Dsicover the children of process 1
	if err := c.walkChildren(0, 1); err != nil {
		return nil, err
	}
	return &c.condb, nil
}
