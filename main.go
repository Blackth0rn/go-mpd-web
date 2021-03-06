package main

import (
	"encoding/json"
	"flag"
	"github.com/fhs/gompd/mpd"
	"github.com/gorilla/websocket"
	"go/build"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

type hub struct {
	connections map[*connection]bool
	inbound     chan []byte
	register    chan *connection
	unregister  chan *connection
	conn        *mpd.Client
	token       int
}

func (h *hub) run() {
	// Configure the connection and handler
	if h.conn == nil {
		conn, err := mpd.Dial("tcp", "localhost:6600")
		if err != nil {
			log.Fatal("mpd.Dial:", err)
		}
		h.conn = conn
	}
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			c.token = h.createConnectionToken()
			jsonReturn, _ := json.Marshal(cmdInput{"register", "", c.token})
			c.send <- jsonReturn
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case msg := <-h.inbound:
			jsonReturn, err, broadcast := h.handleMessage(msg)
			if err != nil {
				jsonReturn, _ = json.Marshal(err)
			}
			for c := range h.connections {
				if broadcast == c.token || broadcast == 0 {
					select {
					case c.send <- jsonReturn:
					default:
						delete(h.connections, c)
						close(c.send)
					}
				}
			}
		case <-time.After(time.Second * 1):
			h.conn.Ping()
		}
	}
}

func (h *hub) handleMessage(m []byte) ([]byte, error, int) {
	return mpdMessageHandle(h.conn, m)
}

func (h *hub) createConnectionToken() int {
	h.token = h.token + 1
	return h.token
}

var h = hub{
	inbound:     make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
	conn:        nil,
	token:       1,
}

type connection struct {
	ws    *websocket.Conn
	send  chan []byte
	token int
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.inbound <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		log.Print(string(message))
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

var (
	addr      = flag.String("addr", ":8080", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/Blackth0rn/go-mpd-web/public", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	http.ServeFile(c, req, filepath.Join(*assets, "index.html"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

func main() {
	flag.Parse()
	go h.run()
	fs := http.FileServer(http.Dir(defaultAssetPath()))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
