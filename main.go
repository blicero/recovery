// /home/krylon/go/src/github.com/blicero/recovery/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 03. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-01 13:56:18 krylon>

package main

import (
	"fmt"
	"os"

	"github.com/blicero/recovery/common"
	"github.com/blicero/recovery/web"
)

func main() {
	fmt.Printf("%s %s (%s) starting up\n",
		common.AppName,
		common.Version,
		common.BuildStamp)

	var (
		srv *web.Server
		err error
	)

	if srv, err = web.Create("localhost:7002", false); err != nil {
		fmt.Fprintf(os.Stderr,
			"Cannot create web server: %s\n",
			err.Error())
		os.Exit(1)
	}

	srv.ListenAndServe()

	fmt.Println("Bye bye")
}
