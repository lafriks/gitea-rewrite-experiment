// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build !bindata

package i18n

import (
	"net/http"
)

var locales http.FileSystem = http.Dir("locale")
