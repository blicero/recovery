// /home/krylon/go/src/github.com/blicero/recovery/database/02_database_mood_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 29. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-01 08:56:27 krylon>

package database

import (
	"math/rand"
	"testing"
	"time"

	"github.com/blicero/recovery/common"
	"github.com/blicero/recovery/data"
)

const (
	minTimestamp = 1648536443 - 608400
	moodCnt      = 28
	timeStep     = 86400 / 4
)

func TestMoodAdd(t *testing.T) {
	if db == nil {
		t.SkipNow()
	}

	var stamp int64 = minTimestamp

	for i := 0; i < moodCnt; i++ {
		var (
			err error
			m   = data.Mood{
				Timestamp: time.Unix(stamp, 0),
				Score:     uint8(rand.Intn(256)),
			}
		)

		if err = db.MoodAdd(&m); err != nil {
			t.Errorf("Cannot add Mood: %s", err.Error())
		}

		stamp += timeStep
	}
} // func TestMoodAdd(t *testing.T)

func TestMoodGetByTime(t *testing.T) {
	if db == nil {
		t.SkipNow()
	}

	var (
		err     error
		results []data.Mood
		t1      = time.Unix(minTimestamp, 0)
		t2      = time.Unix(minTimestamp+moodCnt*timeStep+86400, 0)
	)

	if results, err = db.MoodGetByTime(t1, t2); err != nil {
		t.Fatalf("Error getting mood by time (%s - %s): %s",
			t1.Format(common.TimestampFormat),
			t2.Format(common.TimestampFormat),
			err.Error())
	} else if len(results) != moodCnt {
		t.Fatalf("Unexpected number of results: %d (expected %d)",
			len(results),
			moodCnt)
	}
} // func TestMoodGetByTime(t *testing.T)

func TestMoodGetMostRecent(t *testing.T) {
	if db == nil {
		t.SkipNow()
	}

	var (
		err     error
		results []data.Mood
	)

	if results, err = db.MoodGetMostRecent(moodCnt); err != nil {
		t.Fatalf("Error getting %d most recent mood entries: %s",
			moodCnt,
			err.Error())
	} else if len(results) != moodCnt {
		t.Fatalf("Unexpected number of results: %d (expected %d)",
			len(results),
			moodCnt)
	}
} // func TestMoodGetMostRecent(t *testing.T)
