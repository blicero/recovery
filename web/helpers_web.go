// /home/krylon/go/src/pepper/web/helpers_web.go
// -*- mode: go; coding: utf-8; -*-
// Created on 04. 09. 2019 by Benjamin Walkenhorst
// (c) 2019 Benjamin Walkenhorst
// Time-stamp: <2022-04-07 09:10:29 krylon>
//
// Helper functions for use by the HTTP request handlers

package web

import (
	"encoding/json"

	"github.com/blicero/recovery/data"
	"gonum.org/v1/gonum/stat"
)

func errJSON(msg string) []byte {
	type emsg struct {
		Status  bool
		Message string
	}

	var m = emsg{Message: msg}
	var buf []byte

	buf, _ = json.Marshal(&m)

	return buf
} // func errJSON(msg string) []byte

func linearRegression(values []data.Mood) (float64, float64) {
	const origin = false
	var (
		offset, slope, idx float64
		xs                 = make([]float64, 0, len(values))
		ys                 = make([]float64, 0, len(values))
		weights            []float64
		// minStamp      = time.Now().Unix() - (86400 * 2)
	)

	for _, v := range values {
		//if v.Timestamp.Unix() >= minStamp {
		xs = append(xs, idx) //append(xs, float64(v.Timestamp.Unix()))
		ys = append(ys, float64(v.Score))
		idx++
		// }
	}

	offset, slope = stat.LinearRegression(xs, ys, weights, origin)

	return offset, slope
} // func LinearRegression(data []data.Mood) (float64, float64)

// func getMimeType(path string) (string, error) {
// 	var (
// 		fh      *os.File
// 		err     error
// 		buffer  [512]byte
// 		byteCnt int
// 	)

// 	if fh, err = os.Open(path); err != nil {
// 		return "", err
// 	}

// 	defer fh.Close() // nolint: errcheck

// 	if byteCnt, err = fh.Read(buffer[:]); err != nil {
// 		return "", fmt.Errorf("cannot read from %s: %s",
// 			path,
// 			err.Error())
// 	}

// 	return http.DetectContentType(buffer[:byteCnt]), nil
// }
// func getMimeType(path string) (string, error)

// Local Variables:  //
// compile-command: "go generate && go vet && go build -v -p 16 && gometalinter && go test -v" //
// End: //
