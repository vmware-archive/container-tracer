// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Enyinna Ochulor (VMware) <enyinnaochulor@gmail.com>
 *
 */
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.eng.vmware.com/opensource/tracecruncher-api/api"
	ctx "gitlab.eng.vmware.com/opensource/tracecruncher-api/internal/tracerctx"
	"net/http/httptest"
	"testing"
)

type testRouterStruct struct{}

func (r *testRouterStruct) SetupRouter() *gin.Engine {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	_, eng := gin.CreateTestContext(w)

	return eng
}

// This validates the presence of a router
func TestNewRouterCreated(t *testing.T) {
	api.Router = &testRouterStruct{}
	nodeRouter := NewRouter(&ctx.Tracer{})
	assert.NotNil(t, nodeRouter)
}
