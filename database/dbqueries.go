// /home/krylon/go/src/github.com/blicero/recovery/database/dbqueries.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-01 08:40:45 krylon>

package database

import "github.com/blicero/recovery/database/query"

var dbQueries = map[query.ID]string{
	query.MoodAdd: "INSERT INTO mood (timestamp, score, note) VALUES (?, ?, ?)",
	query.MoodGetByTime: `
SELECT
    id,
    timestamp,
    score,
    note
FROM mood
WHERE timestamp BETWEEN ? AND ?
ORDER BY timestamp
`,
	query.MoodGetMostRecent: `
SELECT
    id,
    timestamp,
    score,
    note
FROM mood
ORDER BY timestamp DESC
LIMIT ?
`,
	query.CravingAdd: "INSERT INTO craving (timestamp, score, note) VALUES (?, ?, ?)",
	query.CravingGetByTime: `
SELECT
    id,
    timestamp,
    score,
    note
FROM craving
WHERE timestamp BETWEEN ? AND ?
ORDER BY timestamp
`,
	query.CravingGetMostRecent: `
SELECT
    id,
    timestamp,
    score,
    note
FROM craving
ORDER BY timestamp DESC
LIMIT ?
`,
}
