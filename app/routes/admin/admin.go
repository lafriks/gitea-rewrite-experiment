// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"code.gitea.io/gitea/models"

	"github.com/gin-gonic/gin"
)

// Dashboard show admin panel dashboard
func (r *Route) Dashboard(c *gin.Context) {
	user := &models.User{}
	if r.DB().First(user).RecordNotFound() {
		user = &models.User{}
		user.Email = "test@example.com"
		if err := r.DB().Save(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		return
	}
	c.HTML(http.StatusOK, "admin/dashboard", gin.H{
		"Version": r.Config().AppVer,
	})
}
