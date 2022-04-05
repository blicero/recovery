// /home/krylon/go/src/github.com/blicero/recovery/database/query/query.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-05 15:00:22 krylon>

//go:generate stringer -type=ID

// Package query contains symbolic constants to represent the various database
// queries we intend to perform.
package query

type ID int8

const (
	MoodAdd ID = iota
	MoodGetByTime
	MoodGetMostRecent
	MoodGetRunningAverage
	CravingAdd
	CravingGetByTime
	CravingGetMostRecent
	CravingGetRunningAverage
)
