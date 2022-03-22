// /home/krylon/go/src/github.com/blicero/recovery/data/data.go
// -*- mode: go; coding: utf-8; -*-
// Created on 16. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-03-22 13:57:35 krylon>

// Package data contains definitions for data types.
package data

import "time"

// Mood represents the user's mood at a given point in time.
type Mood struct {
	ID        int64
	Timestamp time.Time
	Score     uint8
	Note      string
}

// Craving represents the user's craving (duh) for their substance of choice
// at a given point in time.
type Craving struct {
	ID        int64
	Timestamp time.Time
	Score     uint8
	Note      string
}
