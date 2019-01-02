// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package i18n

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/leonelquinteros/gotext"
)

type TranslationStorage struct {
	custom       http.Dir
	locales      map[string]*Locale
	compiled     bool
	disableCache bool

	// Sync Mutex
	sync.RWMutex
}

func New(customDir string, compiled bool) *TranslationStorage {
	return &TranslationStorage{
		custom:   http.Dir(customDir),
		locales:  make(map[string]*Locale),
		compiled: compiled,
	}
}

func NewWithoutCache(customDir string, compiled bool) *TranslationStorage {
	return &TranslationStorage{
		custom:       http.Dir(customDir),
		locales:      make(map[string]*Locale),
		compiled:     compiled,
		disableCache: true,
	}
}

func (tr *TranslationStorage) getContent(file string) ([]byte, error) {
	var f http.File
	var err error
	f, err = tr.custom.Open(file)
	if err != nil {
		f, err = locales.Open(file)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (tr *TranslationStorage) Get(lang string) *Locale {
	var ok bool
	var l *Locale

	lang = gotext.SimplifiedLocale(lang)
	if !tr.disableCache {
		tr.RLock()
		l, ok = tr.locales[lang]
		tr.RUnlock()
	}
	if !ok {
		l = NewLocale(lang)
		var t gotext.Translator
		if tr.compiled {
			if b, err := tr.getContent(lang + ".mo"); err == nil {
				t = &gotext.Mo{}
				t.Parse(b)
			}
		}
		if t == nil {
			if b, err := tr.getContent(lang + ".po"); err == nil {
				t = &gotext.Po{}
				t.Parse(b)
			}
		}
		if t == nil {
			t = &gotext.Po{}
		}
		l.AddDomain("gitea", t)

		if !tr.disableCache {
			tr.Lock()
			tr.locales[lang] = l
			tr.Unlock()
		}
	}

	return l
}
