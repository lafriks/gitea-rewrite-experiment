// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB struct {
	*gorm.DB

	Tables []interface{}
}

func New() (*DB, error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	m := &DB{DB: db, Tables: make([]interface{}, 0)}
	initTables(m)

	return m, nil
}

func initTables(db *DB) {
	db.Tables = append(db.Tables,
		&User{},
	)
}

func (db *DB) Migrate() error {
	// TODO: Custom migrations

	db.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(db.Tables...)
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}
