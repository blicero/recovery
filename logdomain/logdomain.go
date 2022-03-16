// /home/krylon/go/src/github.com/blicero/recovery/logdomain/logdomain.go
// -*- mode: go; coding: utf-8; -*-
// Created on 16. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-03-16 15:49:53 krylon>

// Package logdomain provides constants to identify the different
// "areas" of the application that perform logging.
package logdomain

//go:generate stringer -type=ID

// ID represents an area of concern.
type ID uint8

// These constants identify the various logging domains.
const (
	Common ID = iota
	DBPool
	Database
	Web
)

// AllDomains returns a slice of all the known log sources.
func AllDomains() []ID {
	return []ID{
		Common,
		DBPool,
		Database,
		Web,
	}
} // func AllDomains() []ID
