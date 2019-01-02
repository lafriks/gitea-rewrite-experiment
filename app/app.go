// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package app

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/i18n"

	"github.com/gin-gonic/gin"
)

// App instance
type App struct {
	env Environment

	config *Config

	router *gin.Engine

	db *models.DB

	locales *i18n.TranslationStorage
}

// New creates new application instance
func New(defaultMode string) *App {
	a := &App{
		env: NewEnvironment(defaultMode),
	}
	if a.Env().IsDevelopment() {
		a.locales = i18n.NewWithoutCache("custom/locale", false)
	} else {
		gin.SetMode(gin.ReleaseMode)
		a.locales = i18n.New("custom/locale", true)
	}

	return a
}

// LoadConfig loads application configuration from file
func (a *App) LoadConfig(filePath string) error {
	a.config = &Config{}
	return nil
}

// SetVersion sets application version and built with tags
func (a *App) SetVersion(version, builtWith string) {
	a.config.AppVer = version
	a.config.AppBuiltWith = builtWith
}

// Config returns application configuration
func (a *App) Config() *Config {
	return a.config
}

// SetRouter sets web router
func (a *App) SetRouter(router *gin.Engine) {
	a.router = router
}

func (a *App) InitDB() error {
	var err error
	if a.db, err = models.New(); err != nil {
		return err
	}

	if err = a.db.Migrate(); err != nil {
		return err
	}

	return nil
}

func (a *App) DB() *models.DB {
	return a.db
}

// Start web application
func (a *App) Start() error {
	return a.router.Run(":3000")
}

func (a *App) GetLocale(lang string) *i18n.Locale {
	return a.locales.Get(lang)
}

func (a *App) Env() Environment {
	return a.env
}
