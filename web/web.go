// /home/krylon/go/src/github.com/blicero/recovery/web/web.go
// -*- mode: go; coding: utf-8; -*-
// Created on 01. 04. 2022 by Benjamin Walkenhorst
// (c) 2022 Benjamin Walkenhorst
// Time-stamp: <2022-04-01 10:42:45 krylon>

// Package web provides the web interface to the application.
package web

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/blicero/krylib"
	"github.com/blicero/recovery/common"
	"github.com/blicero/recovery/database"
	"github.com/blicero/recovery/logdomain"
	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
)

//go:embed html
var assets embed.FS

const (
	defaultPoolSize = 4
)

// Server implements the web interface
type Server struct {
	Addr      string
	web       http.Server
	log       *log.Logger
	msgBuf    *krylib.MessageBuffer
	router    *mux.Router
	tmpl      *template.Template
	mimeTypes map[string]string
	pool      *database.Pool
}

// Create creates a new Server instance.
func Create(addr string, keepAlive bool) (*Server, error) {
	var (
		err error
		msg string
		srv = &Server{
			Addr:   addr,
			msgBuf: krylib.CreateMessageBuffer(),
			mimeTypes: map[string]string{
				".css": "text/css",
				".map": "application/json",
				".js":  "text/javascript",
				".png": "image/png",
			},
		}
	)

	if srv.log, err = common.GetLogger(logdomain.Web); err != nil {
		return nil, err
	} else if srv.pool, err = database.NewPool(defaultPoolSize); err != nil {
		srv.log.Printf("[ERROR] Cannot create DB pool: %s\n",
			err.Error())
		return nil, err
	}

	const tmplFolder = "html/templates"
	var templates []fs.DirEntry
	var tmplRe = regexp.MustCompile("[.]tmpl$")

	if templates, err = assets.ReadDir(tmplFolder); err != nil {
		srv.log.Printf("[ERROR] Cannot read embedded templates: %s\n",
			err.Error())
		return nil, err
	}

	srv.tmpl = template.New("").Funcs(funcmap)
	for _, entry := range templates {
		var (
			content []byte
			path    = filepath.Join(tmplFolder, entry.Name())
		)

		if !tmplRe.MatchString(entry.Name()) {
			continue
		} else if content, err = assets.ReadFile(path); err != nil {
			msg = fmt.Sprintf("Cannot read embedded file %s: %s",
				path,
				err.Error())
			srv.log.Printf("[CRITICAL] %s\n", msg)
			return nil, errors.New(msg)
		} else if srv.tmpl, err = srv.tmpl.Parse(string(content)); err != nil {
			msg = fmt.Sprintf("Could not parse template %s: %s",
				entry.Name(),
				err.Error())
			srv.log.Println("[CRITICAL] " + msg)
			return nil, errors.New(msg)
		} else if common.Debug {
			srv.log.Printf("[TRACE] Template \"%s\" was parsed successfully.\n",
				entry.Name())
		}
	}

	srv.router = mux.NewRouter()
	srv.web.Addr = addr
	srv.web.ErrorLog = srv.log
	srv.web.Handler = srv.router

	srv.router.HandleFunc("/favicon.ico", srv.handleFavIco)
	srv.router.HandleFunc("/static/{file}", srv.handleStaticFile)
	srv.router.HandleFunc("/{page:(?i)(?:index|main)?$}", srv.handleIndex)
	// srv.router.HandleFunc("/mood_submit", srv.handleMoodSubmit)

	srv.router.HandleFunc("/ajax/beacon", srv.handleBeacon)
	srv.router.HandleFunc("/ajax/get_messages", srv.handleGetNewMessages)

	if !common.Debug {
		srv.web.SetKeepAlivesEnabled(keepAlive)
	}

	return srv, nil
} // func Create(addr string, keepAlive bool) (*Server, error)

// ListenAndServe enters the HTTP server's main loop, i.e.
// this method must be called for the Web frontend to handle
// requests.
func (srv *Server) ListenAndServe() {
	var err error

	defer srv.log.Println("[INFO] Web server is shutting down")

	srv.log.Printf("[INFO] Web frontend is going online at %s\n", srv.Addr)
	http.Handle("/", srv.router)

	if err = srv.web.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			srv.log.Printf("[ERROR] ListenAndServe returned an error: %s\n",
				err.Error())
		} else {
			srv.log.Println("[INFO] HTTP Server has shut down.")
		}
	}
} // func (srv *Server) ListenAndServe()

// SendMessage send a message to the server's message queue
func (srv *Server) SendMessage(msg string) {
	srv.msgBuf.AddMessage(msg)

	if common.Debug {
		srv.log.Printf("[DEBUG] MessageBuffer holds %d messages\n",
			srv.msgBuf.Count())
	}
} // func (srv *Server) SendMessage(msg string)

// Close shuts down the server.
func (srv *Server) Close() error {
	var err error

	if err = srv.pool.Close(); err != nil {
		srv.log.Printf("[ERROR] Cannot close database pool: %s\n",
			err.Error())
		return err
	} else if err = srv.web.Close(); err != nil {
		srv.log.Printf("[ERROR] Cannot shutdown HTTP server: %s\n",
			err.Error())
		return err
	}

	return nil
} // func (srv *Server) Close() error

// nolint: unused
func (srv *Server) getMessages() []message {
	var m1 = srv.msgBuf.GetAllMessages()
	var m2 = make([]message, len(m1))

	for idx, m := range m1 {
		m2[idx] = message{
			Timestamp: m.Stamp,
			Message:   m.Msg,
			Level:     "DEBUG",
		}
	}

	return m2
} // func (srv *Server) getMessages() []krylib.Message

////////////////////////////////////////////////////////////////////////////////
//// URL handlers //////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////
////////////// Web UI ///////////////////
/////////////////////////////////////////

func (srv *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	srv.log.Printf("[TRACE] Handle request for %s\n",
		r.URL.EscapedPath())

	const (
		tmplName = "index"
		timespan = 86400 * 14
	)

	var (
		err        error
		msg        string
		begin, end time.Time
		db         *database.Database
		tmpl       *template.Template
		data       = tmplDataIndex{
			tmplDataBase: tmplDataBase{
				Title: "Main",
				Debug: common.Debug,
				URL:   r.URL.String(),
			},
		}
	)

	end = time.Now()
	begin = end.Add(-(time.Second * timespan))

	if tmpl = srv.tmpl.Lookup(tmplName); tmpl == nil {
		msg = fmt.Sprintf("Could not find template %q", tmplName)
		srv.log.Println("[CRITICAL] " + msg)
		srv.sendErrorMessage(w, msg)
		return
	}

	db = srv.pool.Get()
	defer srv.pool.Put(db)

	if data.Mood, err = db.MoodGetByTime(begin, end); err != nil {
		msg = fmt.Sprintf("Cannot load mood data: %s",
			err.Error())
		srv.log.Println("[ERROR] " + msg)
		srv.SendMessage(msg)
		srv.sendErrorMessage(w, msg)
		return
	} else if data.Craving, err = db.CravingGetByTime(begin, end); err != nil {
		msg = fmt.Sprintf("Cannot load craving data: %s",
			err.Error())
		srv.log.Println("[ERROR] " + msg)
		srv.SendMessage(msg)
		srv.sendErrorMessage(w, msg)
		return
	}

	data.Messages = srv.getMessages()

	w.Header().Set("Cache-Control", "no-store, max-age=0")
	if err = tmpl.Execute(w, &data); err != nil {
		msg = fmt.Sprintf("Error rendering template %q: %s",
			tmplName,
			err.Error())
		srv.SendMessage(msg)
		srv.sendErrorMessage(w, msg)
	}
} // func (srv *Server) handleIndex(w http.ResponseWriter, r *http.Request)

/////////////////////////////////////////
////////////// Other ////////////////////
/////////////////////////////////////////

func (srv *Server) handleFavIco(w http.ResponseWriter, request *http.Request) {
	srv.log.Printf("[TRACE] Handle request for %s\n",
		request.URL.EscapedPath())

	const (
		filename = "html/static/favicon.ico"
		mimeType = "image/vnd.microsoft.icon"
	)

	w.Header().Set("Content-Type", mimeType)

	if !common.Debug {
		w.Header().Set("Cache-Control", "max-age=7200")
	} else {
		w.Header().Set("Cache-Control", "no-store, max-age=0")
	}

	var (
		err error
		fh  fs.File
	)

	if fh, err = assets.Open(filename); err != nil {
		msg := fmt.Sprintf("ERROR - cannot find file %s", filename)
		srv.sendErrorMessage(w, msg)
	} else {
		defer fh.Close()
		w.WriteHeader(200)
		io.Copy(w, fh) // nolint: errcheck
	}
} // func (srv *Server) handleFavIco(w http.ResponseWriter, request *http.Request)

func (srv *Server) handleStaticFile(w http.ResponseWriter, request *http.Request) {
	srv.log.Printf("[TRACE] Handle request for %s\n",
		request.URL.EscapedPath())

	// Since we controll what static files the server has available, we
	// can easily map MIME type to slice. Soon.

	vars := mux.Vars(request)
	filename := vars["file"]
	path := filepath.Join("html", "static", filename)

	var mimeType string

	srv.log.Printf("[TRACE] Delivering static file %s to client\n", filename)

	var match []string

	if match = common.SuffixPattern.FindStringSubmatch(filename); match == nil {
		mimeType = "text/plain"
	} else if mime, ok := srv.mimeTypes[match[1]]; ok {
		mimeType = mime
	} else {
		srv.log.Printf("[ERROR] Did not find MIME type for %s\n", filename)
	}

	w.Header().Set("Content-Type", mimeType)

	if common.Debug {
		w.Header().Set("Cache-Control", "no-store, max-age=0")
	} else {
		w.Header().Set("Cache-Control", "max-age=7200")
	}

	var (
		err error
		fh  fs.File
	)

	if fh, err = assets.Open(path); err != nil {
		msg := fmt.Sprintf("ERROR - cannot find file %s", path)
		srv.sendErrorMessage(w, msg)
	} else {
		defer fh.Close()
		w.WriteHeader(200)
		io.Copy(w, fh) // nolint: errcheck
	}
} // func (srv *Server) handleStaticFile(w http.ResponseWriter, request *http.Request)

func (srv *Server) sendErrorMessage(w http.ResponseWriter, msg string) {
	html := `
<!DOCTYPE html>
<html>
  <head>
    <title>Internal Error</title>
  </head>
  <body>
    <h1>Internal Error</h1>
    <hr />
    We are sorry to inform you an internal application error has occured:<br />
    %s
    <p>
    Back to <a href="/index">Homepage</a>
    <hr />
    &copy; 2018 <a href="mailto:krylon@gmx.net">Benjamin Walkenhorst</a>
  </body>
</html>
`

	srv.log.Printf("[ERROR] %s\n", msg)

	output := fmt.Sprintf(html, msg)
	w.WriteHeader(500)
	_, _ = w.Write([]byte(output)) // nolint: gosec
} // func (srv *Server) sendErrorMessage(w http.ResponseWriter, msg string)

////////////////////////////////////////////////////////////////////////////////
//// Ajax handlers /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// const success = "Success"

// func (srv *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
// 	srv.log.Printf("[TRACE] Handle request for %s\n",
// 		r.URL.EscapedPath())

// 	// var (
// 	// 	err        error
// 	// 	msg        string
// 	// 	begin, end time.Time
// 	// 	db         *database.Database
// 	// )
// } // func (srv *Server) handleGetData(w http.ResponseWriter, r *http.Request)

func (srv *Server) handleBeacon(w http.ResponseWriter, r *http.Request) {
	// srv.log.Printf("[TRACE] Handle %s from %s\n",
	// 	r.URL,
	// 	r.RemoteAddr)
	var timestamp = time.Now().Format(common.TimestampFormat)
	const appName = common.AppName + " " + common.Version
	var jstr = fmt.Sprintf(`{ "Status": true, "Message": "%s", "Timestamp": "%s", "Hostname": "%s" }`,
		appName,
		timestamp,
		hostname())
	var response = []byte(jstr)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.WriteHeader(200)
	w.Write(response) // nolint: errcheck,gosec
} // func (srv *Web) handleBeacon(w http.ResponseWriter, r *http.Request)

func (srv *Server) handleGetNewMessages(w http.ResponseWriter, r *http.Request) {
	// srv.log.Printf("[TRACE] Handle %s from %s\n",
	// 	r.URL,
	// 	r.RemoteAddr)

	type msgItem struct {
		Time    string
		Level   logutils.LogLevel
		Message string
	}

	type resBody struct {
		Status   bool
		Message  string
		Messages []msgItem
	}

	var messages = srv.getMessages()
	var res = resBody{
		Status:   true,
		Messages: make([]msgItem, len(messages)),
	}

	for idx, i := range messages {
		res.Messages[idx] = msgItem{
			Time:    i.TimeString(),
			Level:   i.Level,
			Message: i.Message,
		}
	}

	var (
		err error
		msg string
		buf []byte
	)

	if buf, err = json.Marshal(&res); err != nil {
		msg = fmt.Sprintf("Error serializing response: %s",
			err.Error())
		srv.SendMessage(msg)
		res.Message = msg
		buf = errJSON(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.WriteHeader(200)
	if _, err = w.Write(buf); err != nil {
		msg = fmt.Sprintf("Failed to send result: %s",
			err.Error())
		srv.log.Println("[ERROR] " + msg)
		srv.SendMessage(msg)
	}
} // func (srv *Server) handleGetNewMessages(w http.ResponseWriter, r *http.Request)
