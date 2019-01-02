// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	"os"
)

type Environment string

func NewEnvironment(defaultMode string) Environment {
	env := os.Getenv("GITEA_MODE")
	if len(env) == 0 {
		env = defaultMode
	}
	if env == "Production" || env == "Staging" {
		return Environment(env)
	}

	return "Development"
}

func (e Environment) IsProduction() bool {
	return e == "Production"
}

func (e Environment) IsStaging() bool {
	return e == "Staging"
}

func (e Environment) IsDevelopment() bool {
	return e == "Development"
}
