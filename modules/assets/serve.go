// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package assets

import (
	"net/http"
	"strings"
)

// FileSystem represents assets filesystem.
type FileSystem struct {
	http.FileSystem
	custom http.FileSystem
}

// Exists checks if file exists
func (fs *FileSystem) Exists(prefix string, filepath string) bool {
	var f http.File
	var err error

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		f, err = fs.custom.Open(p)
		if err != nil {
			f, err = fs.Open(p)
			if err != nil {
				return false
			}
		}

		defer f.Close()
		s, err := f.Stat()
		if err != nil {
			return false
		}
		return !s.IsDir()
	}
	return false
}

func (fs *FileSystem) Open(path string) (http.File, error) {
	f, err := fs.custom.Open(path)
	if err != nil {
		return fs.FileSystem.Open(path)
	}

	return f, err
}

// New assets filesystem.
func New(customDir string) *FileSystem {
	return &FileSystem{
		FileSystem: assets,
		custom:     http.Dir(customDir),
	}
}
