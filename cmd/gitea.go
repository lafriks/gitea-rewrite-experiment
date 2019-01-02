// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Tags holds the build tags used
var Tags string

// DefaultMode hods the default mode
var DefaultMode string

// RootCmd represents the base command when called without any subcommands
var RootCmd *cobra.Command

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd = &cobra.Command{
		Use:   "gitea",
		Short: "A painless self-hosted Git service",
		Long: `By default, gitea will start serving using the webserver with no
	arguments - which can alternatively be run by running the subcommand web.`,
		RunE: runWeb,
	}
}
