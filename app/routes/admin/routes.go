// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"code.gitea.io/gitea/app"

	"github.com/gin-gonic/gin"
)

// Route contains all admin routes
type Route struct {
	*app.App
}

// InitRoutes initializes routes
func InitRoutes(a *app.App, router *gin.RouterGroup) {
	r := &Route{a}

	router.GET("/dashboard", r.Dashboard)
}
