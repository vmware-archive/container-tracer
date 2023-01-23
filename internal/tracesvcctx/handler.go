// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Tzvetomir Stoyanov (VMware) <tz.stoyanov@gmail.com>
 *
 * Frontend handlers of the container-tracer REST API.
 */
package tracekubectx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	kapi "k8s.io/api/core/v1"
)

func (t *TraceKube) sendRequest(pod string, req *http.Request, body io.ReadCloser) (*http.Response, error) {
	var err error
	var pReq *http.Request

	url := fmt.Sprintf("%s%s", t.tracers[pod].target.String(), req.RequestURI)
	if pReq, err = http.NewRequest(req.Method, url, body); err != nil {
		return nil, err
	}
	pReq.Header = req.Header
	return t.tracers[pod].client.Do(pReq)
}

func aggregateMap(agg interface{}, resp *http.Response) interface{} {
	var aggMap map[string]interface{}
	aMap := make(map[string]interface{})

	if agg == nil {
		aggMap = make(map[string]interface{})
	} else {
		aggMap = agg.(map[string]interface{})
	}

	if data, err := io.ReadAll(resp.Body); err == nil {
		if err := json.Unmarshal(data, &aMap); err == nil {
			for indx, item := range aMap {
				aggMap[indx] = item
			}
		}
	}

	return aggMap
}

func (t *TraceKube) proxySend(c *gin.Context, any bool) error {
	var aggregatedData interface{}
	var reqData []byte
	var err error
	errMsg := "Connection to tracers failed"
	status := http.StatusInternalServerError

	if reqData, err = ioutil.ReadAll(c.Request.Body); err != nil {
		return err
	}

	t.trSync.RLock()
	for n, p := range t.tracers {
		if p.state != kapi.PodRunning || p.client == nil {
			continue
		}
		fmt.Print("\tForward a request to ", p.target, " ... ")
		reqBody := ioutil.NopCloser(bytes.NewReader(reqData))
		if resp, err := t.sendRequest(n, c.Request, reqBody); err == nil {
			fmt.Print(resp.StatusCode, "\n")
			if resp.StatusCode == http.StatusOK {
				status = resp.StatusCode
				aggregatedData = aggregateMap(aggregatedData, resp)
				if any == true {
					resp.Body.Close()
					break
				}
			} else if status != http.StatusOK {
				/* Save last received error, if there is still no StatusOK */
				status = resp.StatusCode
				r, _ := io.ReadAll(resp.Body)
				errMsg = string(r)
				fmt.Print("\t\t", errMsg, "\n")
			}
			resp.Body.Close()
		} else {
			fmt.Print(err, "\n")
		}
	}
	t.trSync.RUnlock()

	if status == http.StatusOK {
		if aggregatedData != nil {
			if d, e := json.Marshal(aggregatedData); e == nil {
				c.Writer.Write(d)
			}
		} else {
			c.JSON(http.StatusOK, "{}")
		}
	} else {
		c.JSON(status, errMsg)
	}

	return nil
}

func (t *TraceKube) ProxyAnyMap(c *gin.Context) {
	t.proxySend(c, true)
}

func (t *TraceKube) ProxyAllMap(c *gin.Context) {
	t.proxySend(c, false)
}
