// /home/krylon/go/src/github.com/blicero/recovery/database/01_database_open_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-03-24 10:22:06 krylon>

package database

import (
	"testing"

	"github.com/blicero/recovery/common"
)

func TestDBOpen(t *testing.T) {
	var err error

	if db, err = Open(common.DbPath); err != nil {
		db = nil
		t.Fatalf("Cannot open/initialize database at %s: %s",
			common.DbPath,
			err.Error())
	}
} // func TestDBOpen(t *testing.T)
