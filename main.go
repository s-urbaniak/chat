package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yosssi/ace"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	tpl, err := ace.Load("index", "", &ace.Options{DynamicReload: true})
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, r.Host)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	con, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	raddr := r.RemoteAddr
	go func() {
		for {
			msgType, p, err := con.ReadMessage()
			if err != nil {
				log.Println("error while reading (exiting)", err)
				return
			}

			var buf bytes.Buffer
			buf.WriteString(raddr)
			buf.WriteString(": ")
			buf.Write(p)

			if err := con.WriteMessage(msgType, buf.Bytes()); err != nil {
				log.Println("error while writing (exiting)", err)
				return
			}
		}
	}()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
