// /home/krylon/go/src/github.com/blicero/moodring/web/tmpl_data.go
// -*- mode: go; coding: utf-8; -*-
// Created on 06. 05. 2020 by Benjamin Walkenhorst
// (c) 2020 Benjamin Walkenhorst
// Time-stamp: <2022-04-05 15:53:48 krylon>
//
// This file contains data structures to be passed to HTML templates.

package web

import (
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/blicero/recovery/common"
	"github.com/blicero/recovery/data"

	"github.com/hashicorp/logutils"
)

type message struct {
	Timestamp time.Time
	Level     logutils.LogLevel
	Message   string
}

func (m *message) TimeString() string {
	return m.Timestamp.Format(common.TimestampFormat)
} // func (m *Message) TimeString() string

func (m *message) Checksum() string {
	var str = m.Timestamp.Format(common.TimestampFormat) + "##" +
		string(m.Level) + "##" +
		m.Message

	var hash = sha512.New()
	hash.Write([]byte(str)) // nolint: gosec,errcheck

	var cksum = hash.Sum(nil)
	var ckstr = fmt.Sprintf("%x", cksum)

	return ckstr
} // func (m *message) Checksum() string

// nolint: deadcode,unused
type tmplDataBase struct {
	Title      string
	Messages   []message
	Debug      bool
	TestMsgGen bool
	URL        string
}

// nolint: deadcode,unused
type tmplDataIndex struct {
	tmplDataBase
	Mood       []data.Mood
	MoodAvg    []data.Mood
	Craving    []data.Craving
	CravingAvg []data.Craving
}

// Local Variables:  //
// compile-command: "go generate && go vet && go build -v -p 16 && gometalinter && go test -v" //
// End: //
