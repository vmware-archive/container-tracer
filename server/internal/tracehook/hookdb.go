// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Internal in-memory database with all available trace helper applications.
 */
package tracehook

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	managerPrefix = "manager."
	defaultPath   = "trace-hooks"
)

type TraceHook struct {
	Name        string
	manager     *hookManager
	Description []string `json:"Description"`
}

type Session struct {
	cmd        *exec.Cmd
	cmdOut     []string
	cmdOutLock sync.RWMutex
	cmdErr     []string
	cmdErrLock sync.RWMutex
	cmdWg      sync.WaitGroup
}

type hookManager struct {
	dir     string
	fexec   string
	Tracers map[string]*TraceHook
}

type TraceHooks struct {
	topDir   *string
	managers map[string]*hookManager
}

func (s *Session) GetOutput() (*[]string, *[]string) {
	out := []string{}
	err := []string{}

	s.cmdOutLock.RLock()
	out = append(out, s.cmdOut...)
	s.cmdOutLock.RUnlock()

	s.cmdErrLock.RLock()
	err = append(err, s.cmdErr...)
	s.cmdErrLock.RUnlock()

	return &out, &err
}

func (h *TraceHooks) GetHook(name *string) (*TraceHook, error) {
	for _, a := range h.managers {
		if tr, ok := a.Tracers[*name]; ok {
			return tr, nil
		}
	}
	return nil, fmt.Errorf("Cannot find trace hook %s", *name)
}

func (h *TraceHooks) scanManagers(dir *string) error {
	/* Walk all subdirectories and look for hook managers */
	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if f.IsDir() {
			p := *dir + "/" + f.Name()
			if e := h.scanManagers(&p); e != nil {
				return e
			}
		} else if strings.HasPrefix(f.Name(), managerPrefix) {
			h.managers[*dir] = &hookManager{
				dir:     *dir,
				fexec:   f.Name(),
				Tracers: make(map[string]*TraceHook),
			}
		}
	}

	return nil
}

/* Call the manager to get available trace hooks and description of each of them */
func (h *TraceHooks) scanTraceHooks(dir *string) error {
	all := exec.Command("./"+h.managers[*dir].fexec, "--get-all")
	all.Dir = *dir

	var allOut bytes.Buffer
	all.Stdout = &allOut
	if e := all.Run(); e != nil {
		/* Failed to run this hook manager, skip it */
		return nil
	}

	for _, s := range strings.Fields(allOut.String()) {
		desc := exec.Command("./"+h.managers[*dir].fexec, "--describe", s)
		desc.Dir = *dir

		var descOut bytes.Buffer
		desc.Stdout = &descOut
		if e := desc.Run(); e != nil {
			continue
		}
		dstrip := []string{}
		for _, d := range strings.Split(descOut.String(), "\n") {
			str := strings.TrimSpace(d)
			if len(str) > 0 {
				dstrip = append(dstrip, str)
			}
		}
		h.managers[*dir].Tracers[s] = &TraceHook{
			Name:        s,
			manager:     h.managers[*dir],
			Description: dstrip,
		}
	}

	return nil
}

func readOutput(s *bufio.Scanner, l *sync.RWMutex, b *[]string) {
	for s.Scan() {
		l.Lock()
		*b = append(*b, s.Text())
		l.Unlock()
	}
}

func (h *TraceHooks) Run(th *TraceHook, pids *[]int, parent *[]int, params *[]string, user *string) (*Session, error) {
	var ret Session

	if pids == nil || len(*pids) < 1 {
		return nil, fmt.Errorf("No tasks are provided")
	}
	args := []string{}
	args = append(args, "--run")
	args = append(args, th.Name)

	args = append(args, "--args")
	sargs := "--pid"
	for _, p := range *pids {
		sargs += " " + strconv.Itoa(p)
	}

	if parent != nil {
		sargs += " --parent"
		for _, p := range *parent {
			sargs += " " + strconv.Itoa(p)
		}
	}

	for _, p := range *params {
		sargs += " " + p
	}
	args = append(args, sargs)
	ret.cmd = exec.Command("./"+th.manager.fexec, args...)
	ret.cmd.Dir = th.manager.dir
	stdoutIn, _ := ret.cmd.StdoutPipe()
	stderrIn, _ := ret.cmd.StderrPipe()
	if err := ret.cmd.Start(); err != nil {
		return nil, err
	}

	scannerOut := bufio.NewScanner(stdoutIn)
	scannerErr := bufio.NewScanner(stderrIn)
	ret.cmdWg.Add(2)
	go func() {
		readOutput(scannerOut, &ret.cmdOutLock, &ret.cmdOut)
		ret.cmdWg.Done()
	}()
	go func() {
		readOutput(scannerErr, &ret.cmdErrLock, &ret.cmdErr)
		ret.cmdWg.Done()
	}()
	return &ret, nil
}

func (h *TraceHooks) Stop(s *Session, wait bool) error {
	if err := s.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	if wait {
		s.cmdWg.Wait()
		if err := s.cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (h *TraceHooks) discoverHooks() error {
	/* Reset trace manager database */
	h.managers = make(map[string]*hookManager)

	/* Traverse through all subdirectories looking for files with 'managerPrefix' */
	if e := h.scanManagers(h.topDir); e != nil {
		return e
	}

	/* Walked through all discovered managers and get available trace hooks */
	for d := range h.managers {
		if e := h.scanTraceHooks(&d); e != nil {
			return e
		}
	}
	return nil
}

/* Create a new database with trace hooks in given directory */
func NewTraceHooksDb(path *string) (*TraceHooks, error) {
	db := TraceHooks{
		topDir: path,
	}

	if db.topDir == nil || *db.topDir == "" {
		db.topDir = &defaultPath
	}

	if e := db.discoverHooks(); e != nil {
		return nil, e
	}

	db.ResetAll()

	return &db, nil
}

func (h *TraceHooks) Get() *map[string]*hookManager {
	return &h.managers
}

/* Reset all tracing subsystems */
func (h *TraceHooks) ResetAll() {
	for d, m := range h.managers {
		cmd := exec.Command("./"+m.fexec, "--clear")
		cmd.Dir = d
		cmd.Run()
	}
}
