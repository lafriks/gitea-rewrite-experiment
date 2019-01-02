// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"github.com/jinzhu/gorm"
)

// User represents the object of individual and member of organization.
type User struct {
	gorm.Model
	LowerName string `gorm:"unique;not null"`
	Name      string `gorm:"unique;not null"`
	FullName  string
	// Email is the primary email address (to be used for communication)
	Email            string `gorm:"not null"`
	KeepEmailPrivate bool
	Password         string `gorm:"not null"`

	// MustChangePassword is an attribute that determines if a user
	// is to change his/her password after registration.
	MustChangePassword bool `gorm:"not null;default false"`
}
