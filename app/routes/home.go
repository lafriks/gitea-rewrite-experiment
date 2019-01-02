// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"

	"code.gitea.io/gitea/models"

	"github.com/gin-gonic/gin"
)

// Home render home page
func (r *Route) Home(c *gin.Context) {
	user := &models.User{}
	if r.DB().First(user).RecordNotFound() {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "home", gin.H{
		"aaa": user.Email,
		"loc": r.App.GetLocale("lv_LV"),
	})
}
