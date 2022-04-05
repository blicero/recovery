// /home/krylon/go/src/github.com/blicero/recovery/database/dbqueries.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-05 15:16:27 krylon>

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
	query.MoodGetRunningAverage: `
SELECT
        m.id,
        m.timestamp,
        CAST(ROUND(AVG(score) FILTER (
                                     WHERE timestamp BETWEEN m.timestamp - ?
                                                         AND m.timestamp)
                   OVER (ORDER BY timestamp)) AS INTEGER)
           AS mavg,
        m.note
FROM mood m
ORDER BY timestamp
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
	query.CravingGetRunningAverage: `
SELECT
        c.id,
        c.timestamp,
        CAST(ROUND(AVG(score) FILTER (
                                     WHERE timestamp BETWEEN c.timestamp - ?
                                                         AND c.timestamp)
                   OVER (ORDER BY timestamp)) AS INTEGER)
           AS mavg,
        c.note
FROM craving c
ORDER BY timestamp
LIMIT ?
`,
}
