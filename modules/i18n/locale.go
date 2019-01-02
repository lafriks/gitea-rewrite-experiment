// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package i18n

import (
	"sync"

	"github.com/leonelquinteros/gotext"
)

type Locale struct {
	// Language for this Locale
	lang string

	// List of available Domains for this locale.
	Domains map[string]gotext.Translator

	// First AddDomain is default Domain
	defaultDomain string

	// Sync Mutex
	sync.RWMutex
}

func NewLocale(lang string) *Locale {
	return &Locale{
		lang: lang,
	}
}

// AddDomain creates a new domain for a given locale object.
// If the domain exists, it gets reloaded.
func (l *Locale) AddDomain(dom string, tr gotext.Translator) {
	// Save new domain
	l.Lock()

	if l.Domains == nil {
		l.Domains = make(map[string]gotext.Translator)
	}
	if l.defaultDomain == "" {
		l.defaultDomain = dom
	}
	l.Domains[dom] = tr

	// Unlock "Save new domain"
	l.Unlock()
}

// GetDomain is the domain getter for the package configuration
func (l *Locale) GetDomain() string {
	l.RLock()
	dom := l.defaultDomain
	l.RUnlock()
	return dom
}

// SetDomain sets the name for the domain to be used.
func (l *Locale) SetDomain(dom string) {
	l.Lock()
	l.defaultDomain = dom
	l.Unlock()
}

// Get uses a domain "default" to return the corresponding Translation of a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) Get(str string, vars ...interface{}) string {
	return l.GetD(l.GetDomain(), str, vars...)
}

// GetN retrieves the (N)th plural form of Translation for the given string in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetN(str, plural string, n int, vars ...interface{}) string {
	return l.GetND(l.GetDomain(), str, plural, n, vars...)
}

// GetD returns the corresponding Translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetD(dom, str string, vars ...interface{}) string {
	return l.GetND(dom, str, str, 1, vars...)
}

// GetND retrieves the (N)th plural form of Translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetND(dom, str, plural string, n int, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].GetN(str, plural, n, vars...)
			}
		}
	}

	// Return the same we received by default
	return gotext.Printf(plural, vars...)
}

// GetC uses a domain "default" to return the corresponding Translation of the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetC(str, ctx string, vars ...interface{}) string {
	return l.GetDC(l.GetDomain(), str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of Translation for the given string in the given context in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	return l.GetNDC(l.GetDomain(), str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetDC(dom, str, ctx string, vars ...interface{}) string {
	return l.GetNDC(dom, str, str, 1, ctx, vars...)
}

// GetNDC retrieves the (N)th plural form of Translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNDC(dom, str, plural string, n int, ctx string, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.Domains != nil {
		if _, ok := l.Domains[dom]; ok {
			if l.Domains[dom] != nil {
				return l.Domains[dom].GetNC(str, plural, n, ctx, vars...)
			}
		}
	}

	// Return the same we received by default
	return gotext.Printf(plural, vars...)
}
