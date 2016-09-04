/*
TODO:
Plenty
Add database?
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/polypmer/ghess"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ChessHandler for persistent game
// in handler functions
type ChessHandler struct {
	g    ghess.Board
	init bool
}

// Index page, link to new game
func (h *ChessHandler) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	html := `
<html>
<link href="/css/style.css" rel="stylesheet">
<h1>Ghess Index</h1>
<a href=/new >New Game</a><br>
<a href=/cxt/1>Text Context?</a>
<a href=/board >View Current Game</a><br>
<br><br><br>
<a href="https://github.com/polypmer/go-chess">Source Code</a>
</html>
`
	fmt.Fprintln(w, html)
}

// newGameHandler creates a new Board object
// And redirects to a board instance /board
func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	h.g = ghess.NewBoard()
	h.init = true
	http.Redirect(w, r, "/board", http.StatusSeeOther)
}

// /board/ route, displays board and new move form.
func (h *ChessHandler) showGameHandler(w http.ResponseWriter,
	r *http.Request) {
	if h.init != true {
		// TODO: FAILURE
		http.Redirect(w, r, "/new/", http.StatusSeeOther)
	}

	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}

	t.Execute(w, h.g.Position())
}

func (h *ChessHandler) testContext(w http.ResponseWriter,
	r *http.Request) {
	newCxt := r.URL.Path[1:]
	fmt.Println(newCxt)
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	// 0.0.0.0 won't work across internal ntwk //10.232.44.100
	PORT := ":8080"
	h := new(ChessHandler)
	h.g = ghess.NewBoard() // This means only playin' one game attime
	hub := newHub()
	go hub.run()
	// Server Routes
	http.HandleFunc("/", h.indexHandler)          // link to new game
	http.HandleFunc("/board/", h.showGameHandler) // view
	http.HandleFunc("/new/", h.newGameHandler)    // new board
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.serveWs(hub, w, r)
	}) // websockets
	http.HandleFunc("/cxt/", h.testContext)
	// Handle Static Files
	// TODO: Combine into one function?
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/",
		http.FileServer(http.Dir("static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/",
		http.FileServer(http.Dir("static/img"))))
	// Why can't I just link them all in the same Handle()?
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("Listening on %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}

/*
Websocket structs and functions?!?
*/

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
// Hub struct for websockets
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// newHub returns a pointer to a new Hub
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// This is the json passed from
// the javascript websockets front end
// It's type dictates what kind of broadcast
// I'll be doing
type inCome struct {
	Type        string `json:"type"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Message     string `json:"message"`
}

// outGo struct sends data back to
// chessboardjs with error position
// and message depending on kind of
// broadcast.
type outGo struct {
	Type     string
	Position string
	Message  string
	Error    string
}

// run starts a Hub select switch for
// accepting and disconnecting clients
// and passing on incoming messages.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is an middleman between the websocket connection and the hub.
// Each Client should have a username
// Attached to whatever is read, and used in chat.
// Also, it should have a 'IsPlaying' bool to allow spectators.
// AND MAYBE it should have a black or white... Yikes this gets
// Complicated.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// Hint: just cause
	c.conn.SetReadLimit(maxMessageSize)              // Why?
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // why?
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// readPump, unlike write, is not a goroutine
	// And each client runs an infinite loop,
	// Until the connection closes, or there is an error
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump(g ghess.Board) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	var feedback string // For sending info to client
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// read json from message
			msg := inCome{}
			json.Unmarshal([]byte(message), &msg)

			switch msg.Type {
			case "move":
				mv := &outGo{}
				err = g.ParseStand(msg.Origin,
					msg.Destination)
				fen := g.Position()
				info := g.Stats()
				check, _ := strconv.ParseBool(info["check"])
				checkmate, _ := strconv.ParseBool(
					info["checkmate"])
				if check {
					feedback = "Check!"
				} else if checkmate {
					feedback = "Checkmate!"
				}
				if err != nil {
					mv = &outGo{
						Type:     "move",
						Position: fen,
						Error:    err.Error(),
					}
				} else {
					mv = &outGo{
						Type:     "move",
						Position: fen,
						Error:    feedback,
					}
				}
				feedback = ""
				j, _ := json.Marshal(mv)
				w.Write([]byte(j))
			case "message":
				chat := &outGo{
					Type:    "message",
					Message: msg.Message,
				}
				j, _ := json.Marshal(chat)
				w.Write([]byte(j))
			case "connection":
				// Should this be put elsewhere?
				chat := &outGo{
					Type:    "connection",
					Message: msg.Message,
				}
				j, _ := json.Marshal(chat)
				w.Write([]byte(j))
			}

			// Close the writer
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage,
				[]byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (h *ChessHandler) serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump(h.g)
	// So every time this handler is called
	// the client reads the pump
	client.readPump()
}
