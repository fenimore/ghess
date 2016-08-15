package main

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
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
	http.HandleFunc("/", boardHandler)
	http.ListenAndServe("0.0.0.0:8080", context.ClearHandler(http.DefaultServeMux))
}
