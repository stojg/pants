package network

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	sendBufferSize = 1024 // How many bytes to keep in the send buffer
)

// Server is struct used for a normal web server and also setting up and
// starting a web socket server
type Server struct {
	port        string
	upgrader    websocket.Upgrader
	broadcast   chan []byte
	incoming    chan []byte
	register    chan *connection
	unregister  chan *connection
	connections map[*connection]bool
}

// NewServer returns a new Server
func NewServer(port string) *Server {
	return &Server{
		port:        port,
		broadcast:   make(chan []byte, 256),
		incoming:    make(chan []byte, 256),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}
}

// Start sets up and starts a static http server and a websocket server
func (server *Server) Start() {
	server.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 10,
		WriteBufferSize: 1024 * 10,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	go server.hub()

	http.HandleFunc("/", server.serveStatic)
	http.HandleFunc("/ws", server.serveWs)

	log.Printf("webserver starting on %s", server.port)
	go server.listenAndServe()
}

// Stop closes all active web socket connections
func (server *Server) Stop() {
	for c := range server.connections {
		if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
			log.Printf("disconnect error %s", err)
		} else {
			log.Printf("client %s disconnected", c.ws.RemoteAddr())
		}
		delete(server.connections, c)
		close(c.send)
	}
}

// Incoming returns an array of messages that has been recieved since the last
// time Incoming was called
func (server *Server) Incoming() [][]byte {
	messages := make([][]byte, 0, len(server.incoming))
	for {
		select {
		case msg := <-server.incoming:
			messages = append(messages, msg)
		default:
			return messages
		}
	}
	return messages
}

// Broadcast sends a message to all currently active web socket connections
func (server *Server) Broadcast(msg []byte) {
	server.broadcast <- msg
}

func (server *Server) listenAndServe() {
	if err := http.ListenAndServe(":"+server.port, nil); err != nil {
		log.Fatalf("http.ListenAndServe: %s ", err)
	}
}

func (server *Server) hub() {
	for {
		select {
		case c := <-server.register:
			server.connections[c] = true
			log.Printf("client %s connected", c.ws.RemoteAddr())
		case c := <-server.unregister:
			if _, ok := server.connections[c]; ok {
				delete(server.connections, c)
				close(c.send)
				log.Printf("client %s disconnected", c.ws.RemoteAddr())
			}
		case msg := <-server.broadcast:
			for c := range server.connections {
				select {
				case c.send <- msg:
				default:
					// delete this connection if there is any problem with it
					close(c.send)
					delete(server.connections, c)
				}
			}
		}
	}
}

// serveStatic serves static files over http
func (server *Server) serveStatic(w http.ResponseWriter, r *http.Request) {
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
func (server *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	socket, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{
		send:   make(chan []byte, sendBufferSize),
		ws:     socket,
		server: server,
	}
	server.register <- c

	// pump messages from this connection to the hub
	go c.readPump()

	// pump messages from the hub to this connection
	go c.writePump()

}
