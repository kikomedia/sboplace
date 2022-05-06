// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var addr = flag.String("addr", ":8080", "game server http service address")
var cert = flag.String("cert", "cert.pem", "certificate filename")
var key = flag.String("key", "key.pem", "certificate key filename")
var wait = flag.Duration("shutdown-timeout", time.Second*15, "Duration Shutdowwn in Seconds")

var server http.Server
var gs *GameServer
var store = sessions.NewCookieStore([]byte("someprettyprivatekey")) // :D

func serveHome(w http.ResponseWriter, r *http.Request) {
	//log.Println(" Home : ", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	// Set some session values.
	session.Values["x"] = "1"

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func serveDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("Create archive...")

	dir := TIMELINE_FOLDER + "/"
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		http.Error(w, "File Not found", http.StatusNotFound)
		return
	}

	buff := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buff)

	for _, file := range files {
		log.Println("filetozip: ", file)
		fileToZip, err := os.Open(dir + file.Name())
		if err != nil {
			http.Error(w, "File Not found", http.StatusNotFound)
			return
		}
		defer fileToZip.Close()

		writer, err := zipWriter.Create(file.Name())
		if err != nil {
			http.Error(w, "Error on writer", http.StatusNotFound)
			return
		}

		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			http.Error(w, "Error Copy", http.StatusNoContent)
			return
		}
	}

	zipWriter.Close()
	readSeeker := bytes.NewReader(buff.Bytes())
	http.ServeContent(w, r, "timeline.zip", time.Now(), readSeeker)
}

func serveImages(w http.ResponseWriter, r *http.Request) {
	//log.Println("Request Image...")
	vars := mux.Vars(r)
	image := vars["image"]

	log.Println(r.URL)
	log.Println(image)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "image/"+image)
}

func prepareGameServer() *GameServer {
	gameserver := newGameServer()
	filename, _ := gameserver.gamedata.getNewestFileName()
	if len(filename) > 0 {
		gameserver.gamedata.loadFromFile(filename)
	}
	fmt.Println("Newest Filename: ", filename)
	go gameserver.run()

	return gameserver
}

func prepareShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		log.Println("...server closed.")
		os.Exit(1)
	}()
}

func runHttpServerWithGameServer(game *GameServer) {
	route := mux.NewRouter()
	fs := http.FileServer(http.Dir("./" + TIMELINE_FOLDER + "/"))
	route.PathPrefix("/timeline/").Handler(http.StripPrefix("/timeline/", fs))
	route.HandleFunc("/", serveHome)
	route.HandleFunc("/zip", serveDownload)
	route.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		if session.Values["x"] == "1" {
			serveWs(game, w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	log.Println("Game server started on port ", *addr)

	server := &http.Server{Addr: *addr, Handler: route}

	var err error

	if len(*cert) > 0 && len(*key) > 0 {
		err = server.ListenAndServeTLS(*cert, *key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}
}

func cleanup() {
	log.Println("Shutting Down...")
	gs.gamedata.saveToFile()
	gs.shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()

	server.Shutdown(ctx)
}

func main() {
	// Prepare
	flag.Parse()
	prepareShutdown()
	gs = prepareGameServer()

	// Run
	log.Println("Starting Gameserver...")
	runHttpServerWithGameServer(gs)
}
