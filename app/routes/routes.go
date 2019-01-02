// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"code.gitea.io/gitea/app"
	"code.gitea.io/gitea/app/routes/admin"
	"code.gitea.io/gitea/app/routes/oauth2"
	"code.gitea.io/gitea/app/routes/wellknown"
	"code.gitea.io/gitea/modules/assets"
	"code.gitea.io/gitea/modules/template"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Route contains all routes
type Route struct {
	*app.App
}

// InitRoutes initializes routes
func InitRoutes(a *app.App) *gin.Engine {
	r := &Route{a}

	router := gin.Default()

	if a.Env().IsDevelopment() {
		router.HTMLRender = template.NewWithoutCache("custom/templates")
	} else {
		router.HTMLRender = template.New("custom/templates")
	}

	router.Use(static.Serve("/", assets.New("custom/public")))
	router.GET("/", r.Home)

	wellknown.InitRoutes(a, router.Group("/.well-known"))
	oauth2.InitRoutes(a, router.Group("/oauth2"))
	admin.InitRoutes(a, router.Group("/admin"))

	return router
}
