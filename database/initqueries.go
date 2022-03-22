// /home/krylon/go/src/github.com/blicero/recovery/database/initqueries.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-03-22 14:15:23 krylon>

package database

var initQueries = []string{
	`
CREATE TABLE mood (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER UNIQUE NOT NULL,
    score INTEGER NOT NULL,
    note TEXT,
    CHECK (score BETWEEN 0 AND 255)
)
`,

	"CREATE INDEX mood_time_idx ON mood (timestamp)",

	`
CREATE TABLE craving (
    id INTEGER PRIMARY KEY,
    timestamp INTEGER UNIQUE NOT NULL,
    score INTEGER NOT NULL,
    note TEXT,
    CHECK (score BETWEEN 0 AND 255)
)
`,

	"CREATE INDEX craving_time_idx ON craving (timestamp)",
}
