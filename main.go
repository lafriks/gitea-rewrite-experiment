// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Gitea (git with a cup of tea) is a painless self-hosted Git Service.
package main // import "code.gitea.io/gitea"

import (
	"code.gitea.io/gitea/cmd"
)

// Version holds the current Gitea version
var Version = "2.0.0-dev"

// Tags holds the build tags used
var Tags = ""

// DefaultMode hods the default mode
var DefaultMode = "Development"

func main() {
	cmd.Tags = Tags
	cmd.DefaultMode = DefaultMode
	cmd.RootCmd.Version = Version
	cmd.Execute()
}
