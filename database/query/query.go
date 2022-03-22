// /home/krylon/go/src/github.com/blicero/recovery/database/query/query.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-03-22 15:08:16 krylon>

//go:generate stringer -type=ID

// Package query contains symbolic constants to represent the various database
// queries we intend to perform.
package query

type ID int8

const (
	MoodAdd ID = iota
	MoodGetByTime
	CravingAdd
	CravingGetByTime
)
