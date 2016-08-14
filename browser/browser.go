package main

import (
	"fmt"
	"github.com/polypmer/ghess"
	"net/http"
)

func boardHandler(g *Board, w http.ResponseWriter, r *http.Request) {
	game := ghess.NewBoard()
	fmt.Fprintf(w, g.String())
}

func main() {
	// Listen
	http.HandleFunc("/", boardHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
