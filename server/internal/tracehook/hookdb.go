// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Internal in-memory database with all available trace helper applications.
 */
package tracehook

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

type hookManager struct {
	fexec   string
	Tracers map[string]*TraceHook
}

type TraceHooks struct {
	topDir   *string
	managers map[string]*hookManager
}

func (c *TraceHooks) GetHook(name *string) (*TraceHook, error) {
	for _, h := range c.managers {
		if tr, ok := h.Tracers[*name]; ok {
			return tr, nil
		}
	}
	return nil, fmt.Errorf("Cannot find trace hook %s", name)
}

func (c *TraceHooks) scanManagers(dir *string) error {
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
			if e := c.scanManagers(&p); e != nil {
				return e
			}
		} else if strings.HasPrefix(f.Name(), managerPrefix) {
			c.managers[*dir] = &hookManager{
				fexec:   f.Name(),
				Tracers: make(map[string]*TraceHook),
			}
		}
	}

	return nil
}

/* Call the manager to get available trace hooks and description of each of them */
func (c *TraceHooks) scanTraceHooks(dir *string) error {
	p, e := os.Getwd()
	if e != nil {
		return e
	}

	if e = os.Chdir(*dir); e != nil {
		return e
	}

	all := exec.Command("./"+c.managers[*dir].fexec, "--get-all")
	var allOut bytes.Buffer
	all.Stdout = &allOut
	if e = all.Run(); e != nil {
		return nil
	}

	for _, s := range strings.Fields(allOut.String()) {
		desc := exec.Command("./"+c.managers[*dir].fexec, "--describe", s)
		var descOut bytes.Buffer
		desc.Stdout = &descOut
		if e = desc.Run(); e != nil {
			continue
		}
		dstrip := []string{}
		for _, d := range strings.Split(descOut.String(), "\n") {
			str := strings.TrimSpace(d)
			if len(str) > 0 {
				dstrip = append(dstrip, str)
			}
		}
		c.managers[*dir].Tracers[s] = &TraceHook{
			Name:        s,
			manager:     c.managers[*dir],
			Description: dstrip,
		}
	}

	if e = os.Chdir(p); e != nil {
		return e
	}
	return nil
}

func (c *TraceHooks) tracerHooksDiscover() error {
	/* Reset trace manager database */
	c.managers = make(map[string]*hookManager)

	/* Traverse through all subdirectories looking for files with 'managerPrefix' */
	if e := c.scanManagers(c.topDir); e != nil {
		return e
	}

	/* Walked through all discovered managers and get available trace hooks */
	for d := range c.managers {
		if e := c.scanTraceHooks(&d); e != nil {
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

	if e := db.tracerHooksDiscover(); e != nil {
		return nil, e
	}

	db.ResetAll()

	return &db, nil
}

func (c *TraceHooks) Get() *map[string]*hookManager {
	return &c.managers
}

/* Reset all tracing subsystems */
func (h *TraceHooks) ResetAll() {
	for d, m := range h.managers {
		cmd := exec.Command("./"+m.fexec, "--clear")
		cmd.Dir = d
		cmd.Run()
	}
}
