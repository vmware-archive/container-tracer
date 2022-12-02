// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2022 VMware, Inc. Enyinna Ochulor (VMware) <enyinnaochulor@gmail.com>
 *
 */
package api

import (
	"github.com/gin-gonic/gin"
)

var (
	Router routerInterface
)

type routerInterface interface {
	SetupRouter() *gin.Engine
}

type routerStruct struct{}

func (r *routerStruct) SetupRouter() *gin.Engine {
	router := gin.Default()

	return router
}

func init() {
	Router = &routerStruct{}
}
