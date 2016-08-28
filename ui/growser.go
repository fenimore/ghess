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
	"time"
)

// ChessHandler for persistant game
// in handler functions
type ChessHandler struct {
	g    ghess.Board
	init bool
}

// Chessboard struct is for
// sending chessboard.js info
type ChessBoard struct {
	Board    string
	Fen      string
	Pgn      string
	Move     string
	wToMove  bool
	Feedback string
}

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

// boardHandler for playing game
// Takes url param pgn move
func (h *ChessHandler) playGameHandler(w http.ResponseWriter,
	r *http.Request) {
	// If no board, redirect to board
	// How to check if struct is empty?
	// Can't compare structs with []byte field?
	// Yuck, TODO
	move := ""
	feedback := ""
	// POST will be for chats?
	if r.Method == "POST" {
		r.ParseForm()
		move = r.Form.Get("move")
	} else if r.Method == "GET" {
		move = r.URL.Path[len("/play/"):]
	}
	e := h.g.ParseMove(move)
	if e != nil {
		feedback = e.Error()
	}
	// Display with GUI chessboard.js
	// TODO: Get White and Black
	// And rotate board accordingly
	b := ChessBoard{Board: h.g.String(), Fen: h.g.Position(), Pgn: h.g.PgnString(), Feedback: feedback}
	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}
	t.Execute(w, b)
}

// newGameHandler creates a new Board object
// and links to /board/ route
func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	h.g = ghess.NewBoard()
	h.init = true
	fmt.Fprintln(w, "<a href=/board>New Game Created</a>")
}

// /board/ route, displays board and new move form.
func (h *ChessHandler) showGameHandler(w http.ResponseWriter,
	r *http.Request) {
	if h.init != true {
		// TODO: FAILURE
		http.Redirect(w, r, "/new/", http.StatusSeeOther)
	}
	b := ChessBoard{Board: h.g.String(), Fen: h.g.Position(), Pgn: h.g.PgnString()}
	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}
	t.Execute(w, b)
}

// Index page, link to new game
func (h *ChessHandler) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintln(w, "<a href=/new >New Game</a>")
}

// AJAX Handler for Updating board
// This does not update all open connections
// TODO: Websockets!?
func (h *ChessHandler) makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	// Get Form Value
	field := r.FormValue("move")
	// make Move
	h.g.ParseMove(field)
	pos := h.g.Position()
	// Write to Client
	w.Write([]byte(pos))
	// TODO: write to all open connections
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	PORT := "0.0.0.0:8080"
	h := new(ChessHandler)
	hub := newHub()
	go hub.run()
	// Server Routes
	http.HandleFunc("/", h.indexHandler)            // link to new game
	http.HandleFunc("/play/", h.playGameHandler)    // deprecated
	http.HandleFunc("/board/", h.showGameHandler)   // view
	http.HandleFunc("/new/", h.newGameHandler)      // new board
	http.HandleFunc("/makemove", h.makeMoveHandler) //ajax
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.serveWs(hub, w, r)
	}) // websockets
	// Handle Static Files
	// TODO: Combine into one function?
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("static/img"))))
	// Why can't I just link them all in the same Handle()?
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//Listen and Server on PORT 8080
	fmt.Printf("Listening on %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}

func getPanel(m map[string]string) string {
	return "|Move:  " + m["move"] + "     Turn: " + m["turn"] +
		"\n|Check: " + m["check"] + " Castle: " + m["castling"] +
		"\n|Mate:  " + m["checkmate"] + " Score: " + m["score"] + "\n"
}

/*
Websocket structs and functions?!?
*/
// hub maintains the set of active clients and broadcasts messages to the
// clients.

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

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

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

			if msg.Type == "move" {
				mv := &outGo{}
				err = g.ParseStand(msg.Origin,
					msg.Destination)
				fen := g.Position()
				if err != nil {
					fmt.Println(err)
					mv = &outGo{
						Type:     "move",
						Position: fen,
						Error:    err.Error(),
					}
				} else {
					mv = &outGo{
						Type:     "move",
						Position: fen,
					}
				}

				j, _ := json.Marshal(mv)
				w.Write([]byte(string(j))) // unnecessary?
			} else if msg.Type == "message" {
				fmt.Println(msg.Message)
			}

			// Close the writer
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
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
	// So everytime this handler is called
	// the client reads the pump
	client.readPump()
}
