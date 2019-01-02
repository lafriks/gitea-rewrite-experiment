// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build !bindata

package template

import (
	"net/http"
)

var templates http.FileSystem = http.Dir("templates")
