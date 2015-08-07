package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type webserver struct {
	port     string
	upgrader websocket.Upgrader
}

func (server *webserver) Start() {

	server.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 10,
		WriteBufferSize: 1024 * 10,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/", server.serveStatic)
	http.HandleFunc("/ws", server.serveWs)

	log.Printf("webserver starting on %s", server.port)
	if err := http.ListenAndServe(":"+server.port, nil); err != nil {
		log.Fatalf("http.ListenAndServe: %s ", err)
	}
}

// serveStatic serves static files over http
func (server *webserver) serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.URL.Path[1:] == "" {
		http.ServeFile(w, r, "public/index.html")
		return
	}
	http.ServeFile(w, r, "public/"+r.URL.Path[1:])
}

// serveWS handles websocket connections
func (server *webserver) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	socket, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: socket}
	h.register <- c

	go c.writePump()
	go c.readPump()
}
