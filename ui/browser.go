package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // only av to sql pkg
	"github.com/polypmer/ghess"
	"net/http"
)

type ChessHandler struct {
	db *sql.DB
}

// boardHandler for playing game
// Takes url param pgn move
func (h *ChessHandler) playGameHandler(w http.ResponseWriter,
	r *http.Request) {
	game := ghess.NewBoard()
	fen := ReadGame(h.db, 1)
	game.LoadFen(fen)
	fmt.Fprintln(w, game.String())
	move := r.URL.Path[1:]
	e := game.ParseMove(move)
	fen = game.Position()
	UpdateRow(h.db, fen, 1)
	if e != nil {
		fmt.Fprintln(w, e.Error())
	}
	fmt.Fprintln(w, game.String())
}

func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	CreateGame(h.db)
	// Create New Game
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	// INIT and CREATE DATABASE
	db := InitDb("./chess.db")
	CreateTable(db)
	//CreateGame(db)
	// TEST EDIT DATABASE
	//ReadGame(db, 1)
	//UpdateRow(db, "LKJDF", 1)
	//ReadGame(db, 1)
	h := ChessHandler{db: db}

	// Server Part
	http.HandleFunc("/", h.playGameHandler)
	http.HandleFunc("/new", h.newGameHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

/* Database Helpers */
// InitDb Open() a sqlite3 database
// with path passed in
func InitDb(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Printf("Error %s in Initiating", err)
	}
	return db
}

// CreateTable in the sql.DB
// Table includes names of players
// The FEN and PGN and Date
func CreateTable(db *sql.DB) {
	sql_table := `
CREATE TABLE IF NOT EXISTS games(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    white TEXT NOT NULL,                
    black TEXT NOT NULL,
    fen TEXT NOT NULL DEFAULT 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1',
    pgn TEXT NOT NULL DEFAULT '',
    date DATE
);
`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Printf("Error %s in Creating\n", err)
	}
}

// Create game
// TODO: get date from time pkg
// TODO: pass in players
func CreateGame(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO games(white, black," +
		"date)values(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec("w", "b", "1989-01-01")
	if err != nil {
		fmt.Println(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("New game: %i\n", id)
}

// Update FEN position of id with fen string
// Change for PGN
func UpdateRow(db *sql.DB, fen string, gId int) {
	stmt, err := db.Prepare("UPDATE GAMES SET fen=? where id=?")
	if err != nil {
		fmt.Printf("Error %s in Update", err)
	}
	// Unlike Query, Exec can
	// Throw away it's result
	_, err = stmt.Exec(fen, gId)
	if err != nil {
		fmt.Printf("Error %s in Result", err)
	}

}

// Read a game give it's id
// Print out the player names and
// The fen position
func ReadGame(db *sql.DB, gId int) string {
	var (
		id    int
		white string
		black string
		fen   string
	)
	rows, err := db.Query("select id, white, black, fen from games where id = ?", gId)
	if err != nil {
		fmt.Printf("Error %s in Query", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &white, &black, &fen)
		if err != nil {
			fmt.Printf("Error %s in Rows", err)
		}
		fmt.Println(id, white, black, fen)
	}
	rows.Close()
	return fen
}

// Read Games
func ReadGames(db *sql.DB) {
	stmt := "SELECT * FROM games"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Print("Error %s in Reading Games", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var white string
		var black string
		var fen string
		var pgn string
		var date string
		err = rows.Scan(&id, &white, &black, &fen, &pgn, &date)
		if err != nil {
			fmt.Printf("Error %s in Scanning rows", err)
		}
		fmt.Println(id)
		fmt.Println(white, black)
	}
	// Redundant but pas mal
	rows.Close()
}
