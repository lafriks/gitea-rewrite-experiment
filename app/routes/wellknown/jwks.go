// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package wellknown

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"
)

// JWKsJSON returns JWKs public keys
func (r *Route) JWKsJSON(c *gin.Context) {
	var jwks jose.JSONWebKeySet

	c.JSON(http.StatusOK, &jwks)
}
