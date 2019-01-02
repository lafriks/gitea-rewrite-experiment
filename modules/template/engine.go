// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var (
	htmlContentType = []string{"text/html; charset=utf-8"}
)

type templateEngine struct {
	custom   http.FileSystem
	config   templateConfig
	tplMap   map[string]*template.Template
	tplMutex sync.RWMutex
}

type templateRender struct {
	Engine *templateEngine
	Name   string
	Data   interface{}
}

type templateConfig struct {
	Extension    string           //template extension
	Master       string           //template master
	Partials     []string         //template partial, such as head, foot
	Funcs        template.FuncMap //template functions
	DisableCache bool             //disable cache, debug mode
	Delims       render.Delims    //delimeters
}

func getConfig(disableCache bool) templateConfig {
	return templateConfig{
		Extension:    ".tmpl",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        NewFuncMap(),
		DisableCache: disableCache,
		Delims:       render.Delims{Left: "{{", Right: "}}"},
	}
}

// New HTML template renderer
func New(customDir string) render.HTMLRender {
	return &templateEngine{
		custom:   http.Dir(customDir),
		config:   getConfig(false),
		tplMap:   make(map[string]*template.Template),
		tplMutex: sync.RWMutex{},
	}
}

// NewWithoutCache HTML template renderer without caching
func NewWithoutCache(customDir string) render.HTMLRender {
	return &templateEngine{
		custom:   http.Dir(customDir),
		config:   getConfig(true),
		tplMap:   make(map[string]*template.Template),
		tplMutex: sync.RWMutex{},
	}
}

func (e *templateEngine) Instance(name string, data interface{}) render.Render {
	return templateRender{
		Engine: e,
		Name:   name,
		Data:   data,
	}
}

func (e *templateEngine) HTML(ctx *gin.Context, code int, name string, data interface{}) {
	instance := e.Instance(name, data)
	ctx.Render(code, instance)
}

func (e *templateEngine) executeRender(out io.Writer, name string, data interface{}) error {
	master := e.config.Master
	if filepath.Ext(name) == e.config.Extension {
		master = ""
		name = strings.TrimSuffix(name, e.config.Extension)

	}
	return e.executeTemplate(out, name, data, master)
}

func (e *templateEngine) executeTemplate(out io.Writer, name string, data interface{}, master string) error {
	var tpl *template.Template
	var err error
	var ok bool

	allFuncs := make(template.FuncMap, 0)
	allFuncs["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := e.executeTemplate(buf, layout, data, master)
		return template.HTML(buf.String()), err
	}

	// Get the plugin collection
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	if !e.config.DisableCache {
		e.tplMutex.RLock()
		tpl, ok = e.tplMap[name]
		e.tplMutex.RUnlock()
	}

	exeName := name
	if len(master) != 0 {
		exeName = master
	}

	if !ok {
		tplList := make([]string, 0)
		if len(master) != 0 {
			tplList = append(tplList, master)
		}
		tplList = append(tplList, name)
		tplList = append(tplList, e.config.Partials...)

		// Loop through each template and test if template exists and parses
		tpl = template.New(name).Funcs(allFuncs).Delims(e.config.Delims.Left, e.config.Delims.Right)
		for _, v := range tplList {
			var data string
			data, err = e.readFileContent(v)
			if err != nil {
				return err
			}
			var tmpl *template.Template
			if v == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(v)
			}
			_, err = tmpl.Parse(data)
			if err != nil {
				return fmt.Errorf("Renderer parser '%v' error: %v", v, err)
			}
		}
		if !e.config.DisableCache {
			e.tplMutex.Lock()
			e.tplMap[name] = tpl
			e.tplMutex.Unlock()
		}
	}

	// Display the content to the screen
	err = tpl.Funcs(allFuncs).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("Error executing template: %v", err)
	}

	return nil
}

func (e *templateEngine) readFileContent(tplFile string) (string, error) {
	var f http.File
	var err error
	f, err = e.custom.Open(tplFile + e.config.Extension)
	if err != nil {
		f, err = templates.Open(tplFile + e.config.Extension)
		if err != nil {
			return "", err
		}
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r templateRender) Render(w http.ResponseWriter) error {
	return r.Engine.executeRender(w, r.Name, r.Data)
}

func (r templateRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
}
