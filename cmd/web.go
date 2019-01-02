// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/gitea/app"
	"code.gitea.io/gitea/app/routes"

	"github.com/spf13/cobra"
)

var port int16
var configPath, pidPath, mode string

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start Gitea web server",
	Long: `Gitea web server is the only thing you need to run,
and it takes care of all the other things for you`,
	RunE: runWeb,
}

func runWeb(cmd *cobra.Command, args []string) error {
	if len(mode) == 0 {
		mode = DefaultMode
	}
	a := app.New(mode)
	if err := a.LoadConfig(configPath); err != nil {
		return err
	}
	a.SetVersion(RootCmd.Version, Tags)
	a.SetRouter(routes.InitRoutes(a))
	if err := a.InitDB(); err != nil {
		return err
	}

	return a.Start()
}

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVarP(&mode,
		"mode",
		"m",
		"",
		"Mode (Development, Staging, Production)")

	webCmd.Flags().Int16VarP(&port,
		"port",
		"p",
		3000,
		"Temporary port number to prevent conflict")

	webCmd.Flags().StringVarP(&configPath,
		"config",
		"c",
		"custom/conf/app.ini",
		"Configuration file path")
	webCmd.Flags().SetAnnotation("config",
		cobra.BashCompFilenameExt,
		[]string{"json", "yaml", "yml", "toml", "hcl"})

	webCmd.Flags().StringVarP(&configPath,
		"pid",
		"P",
		"/var/run/gitea.pid",
		"pid file path")
	webCmd.Flags().SetAnnotation("pid",
		cobra.BashCompFilenameExt,
		[]string{"pid"})
}
