// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"runtime"
	"strings"
)

// NewFuncMap returns functions for injecting to templates
func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"GoVer": func() string {
			return strings.Title(runtime.Version())
		},
		"Safe":   Safe,
		"SafeJS": SafeJS,
		"Add": func(a, b int) int {
			return a + b
		},
		"List": List,
		"SubStr": func(str string, start, length int) string {
			if len(str) == 0 {
				return ""
			}
			end := start + length
			if length == -1 {
				end = len(str)
			}
			if len(str) < end {
				return str
			}
			return str[start:end]
		},
		"PathEscape": url.PathEscape,
		"JsonPrettyPrint": func(in string) string {
			var out bytes.Buffer
			err := json.Indent(&out, []byte(in), "", "  ")
			if err != nil {
				return ""
			}
			return out.String()
		},
		"Printf": fmt.Sprintf,
		"Escape": Escape,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values) == 0 {
				return nil, errors.New("invalid dict call")
			}

			dict := make(map[string]interface{})

			for i := 0; i < len(values); i++ {
				switch key := values[i].(type) {
				case string:
					i++
					if i == len(values) {
						return nil, errors.New("specify the key for non array values")
					}
					dict[key] = values[i]
				case map[string]interface{}:
					m := values[i].(map[string]interface{})
					for i, v := range m {
						dict[i] = v
					}
				default:
					return nil, errors.New("dict values must be maps")
				}
			}
			return dict, nil
		},
	}
}

// Safe render raw as HTML
func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

// SafeJS renders raw as JS
func SafeJS(raw string) template.JS {
	return template.JS(raw)
}

// Escape escapes a HTML string
func Escape(raw string) string {
	return html.EscapeString(raw)
}

// List traversings the list
func List(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}
