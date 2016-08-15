package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // only av to sql pkg
	"github.com/polypmer/ghess"
	"net/http"
)

func boardHandler(w http.ResponseWriter, r *http.Request) {
	game := ghess.NewBoard()
	fmt.Fprintln(w, game.String())
	move := r.URL.Path[1:]
	e := game.ParseMove(move)
	if e != nil {
		fmt.Fprintln(w, e.Error())
	}
	fmt.Fprintln(w, game.String())
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	db := InitDb("./chess.db")
	CreateTable(db)
	CreateGame(db)

	http.HandleFunc("/", boardHandler)
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
    white TEXT,                
    black TEXT,
    fen TEXT NULL,
    pgn TEXT NULL,
    date DATE
);
`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Printf("Error %s in Creating\n", err)
	}
}

// Create game
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
	fmt.Printf("New game: %s\n", string(id))
}

// Update/ Fen
