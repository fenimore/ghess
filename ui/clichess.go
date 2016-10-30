package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/polypmer/ghess"
)

var welcome string = `
=================
	 go-chess

    /~ |_ _  _ _
    \_|||(/__\_\


    Enter /help for more options

`
var manuel string = `
    ====================
	/~ |_ _  _ _
	\_|||(/__\_\
    ====================

go-chess command line client
Fenimore Love 2016

Help:
    Prefix commands with /  (slash)
    C-+ to enlarge font size

PGN Input Example:
    e4
    e5
    Nf3

Commands:
    quit or exit - exit game
    new          - new game
    print        - print game
    panel        - print game info
    coordinates  - print board coordinates
    pgn          - print PGN history
    fen          - print FEN position
    score        - print win/loss
    set-headers  - set PGN headers
    headers      - print game info
    random-game  - play random game
    rand         - make a random move
    valid        - show valid moves
    eval         - score position (+ White | - Black)
    minimax      - generate minimax score
    aivsai       - watch AI match
    aivshuman    - play against AI (human is black)
    aivsrand     - ai vs random

`
var turn string

func main() {
	game := ghess.NewBoard()
	PlayGame(game)
}

// PlayGame is the command line user interface.
// It is used mostly for debugging.
func PlayGame(game ghess.Board) { // TODO Rotate Board
	var info map[string]string

	reader := bufio.NewReader(os.Stdin)
	// welcome message
	fmt.Println(welcome)
	fmt.Print(game.String())
Loop:
	for {
		info = game.Stats()
		if info["turn"] == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Print(turn, " to move: ")
		input, _ := reader.ReadString('\n')
		isCmd, _ := regexp.MatchString(`/`, input)
		if isCmd {
			input = strings.TrimRight(input, "\r\n")
			switch input {
			case "/help":
				fmt.Print("\n", manuel)
			case "/quit":
				break Loop
			case "/exit":
				break Loop
			case "/new":
				game = ghess.NewBoard()
				fmt.Print(game.String())
			case "/print":
				fmt.Print(game.String())
			case "/panel":
				info = game.Stats()
				fmt.Print(getPanel(info))
			case "/coordinates":
				fmt.Println("Coordinates:")
				game.Coordinates()
			case "/score":
				checkMate, _ := strconv.ParseBool(info["checkmate"])
				score := info["score"]
				if checkMate { // TODO: or draw
					fmt.Println("Game over")
				}
				fmt.Println("Score: ", score)
			case "/pgn":
				fmt.Println("PGN history:")
				fmt.Println(game.PgnString())
			case "/load-pgn":
				var err error
				fmt.Print("Enter PGN history: ")
				history, _ := reader.ReadString('\n')
				err = game.LoadPgn(history)
				if err != nil {
					fmt.Println(err)
				}
				info := game.Stats()
				fmt.Print(getPanel(info))
				fmt.Print(game.String())
			case "/load-fen":
				var err error
				fmt.Print("Enter FEN position: ")
				position, _ := reader.ReadString('\n')
				err = game.LoadFen(position)
				if err != nil {
					fmt.Println(err)
				}
				info := game.Stats()
				fmt.Print(getPanel(info))
				fmt.Print(game.String())
			case "/fen":
				fmt.Println("FEN position:")
				fmt.Println(game.Position())
			case "/set-headers":
				fmt.Print("Enter White Player: ")
				first, _ := reader.ReadString('\n')
				fmt.Print("Enter Black Player: ")
				second, _ := reader.ReadString('\n')
				game.SetHeaders(first, second)
			case "/headers":
				info = game.Stats()
				fmt.Println(info["headers"])
			case "/valid":
				origs, dests := game.SearchValid()
				fmt.Println(origs)
				fmt.Println(dests)
				fmt.Println("Total valid moves: ",
					len(origs))
			case "/eval":
				score := game.Evaluate()
				fmt.Println("Position: ", score)
			case "/computer":
				// Play as white against the computer
			case "/ai":
				makeAiMove(&game)
			case "/prune":
				s, e := ghess.MiniMaxPruning(0, 4, ghess.GetState(&game))
				fmt.Println(s, e)
			case "/minimax":
				predictAiMove(game)
			case "/aivsai":
				aiVsAi(game)
			case "/aivsrand":
				aiVsRandom(game)
			case "/aivshuman":
				aiVsHuman(game, reader)
			case "/rand":
				origs, dests := game.SearchValid()
				e := game.MoveRandom(origs, dests)
				if e != nil {
					fmt.Println(e)
				}
			case "/random-game":
				randomGame(game)
			case "/tension":
				fmt.Println("Tension Coordinates:")
				fmt.Println(game.Tension())
				fmt.Println("Total Tension: ", game.TensionSum())
				fmt.Println(game.StringTension())
			default:
				fmt.Println("Mysterious input")
			}
			continue
		}
		makeMove(&game, input)
	}
	fmt.Println("\nGood Game.")
}

func makeMove(game *ghess.Board, input string) {
	e := game.ParseMove(input)
	info := game.Stats()
	if info["turn"] == "w" {
		turn = "White"
	} else {
		turn = "Black"
	}
	fmt.Println("\n_________________")
	info = game.Stats()
	fmt.Print(getPanel(info))
	// TODO use formats.
	if e != nil {
		fmt.Printf("|   [Error: %v]\n", e)
	}
	fmt.Print(game.StringWhite())
	ch, _ := strconv.ParseBool(info["check"])
	checkmate, _ := strconv.ParseBool(info["checkmate"])
	if checkmate {
		fmt.Println("****Check and Mate.****")
	} else if ch {
		fmt.Println("****Check!****")
	}

}

func randomGame(game ghess.Board) {
	reader := bufio.NewReader(os.Stdin)
	info := game.Stats()
	exit := false
	go func() {
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		// Todo turn into chan
		exit = true
	}()
	fmt.Println("\nPress Return to stop")
	time.Sleep(2000 * time.Millisecond)
LoopRand:
	for {
		if exit == true {
			break LoopRand
		}
		origs, dests := game.SearchValid()
		game.MoveRandom(origs, dests)
		info := game.Stats()
		fmt.Println("Move ", info["move"])
		fmt.Print(game.StringWhite())
		think(true)
		gameOver, _ := strconv.ParseBool(info["checkmate"])
		check, _ := strconv.ParseBool(info["check"])
		if check {
			fmt.Println("****Check****")
		}
		if gameOver {
			break
		}
	}
	fmt.Println(info["score"])

}

func aiVsRandom(game ghess.Board) {
VsLoop:
	for {
		state, err := ghess.MiniMax(0, 3, ghess.GetState(&game))
		if err != nil {
			fmt.Println(err)
			break VsLoop
		}
		game.Move(state.Init[0], state.Init[1])
		info := game.Stats()
		fmt.Println(getPanel(info))
		fmt.Println(game.StringWhite())
		if info["score"] != "*" {
			fmt.Println("Game Over:")
			fmt.Println(info["score"])
			break VsLoop
		}

		origs, dests := game.SearchValid()
		e := game.MoveRandom(origs, dests)
		think(true)
		if e != nil {
			fmt.Println(e)
		}
		info = game.Stats()
		fmt.Println(getPanel(info))
		fmt.Println(game.StringWhite())
		if info["score"] != "*" {
			fmt.Println("Game Over:")
			fmt.Println(info["score"])
			break VsLoop
		}
	}

}

func aiVsAi(game ghess.Board) {
	// TODO, AI versus weaker AI
AiLoop:
	for {
		now := time.Now()
		state, err := ghess.MiniMax(0, 3, ghess.GetState(&game))
		if err != nil {
			fmt.Println(err)
			break AiLoop
		}
		game.Move(state.Init[0], state.Init[1])
		info := game.Stats()
		fmt.Println(getPanel(info))
		fmt.Println(game.StringWhite())
		fmt.Printf("\nThis took me: %s\n", time.Since(now))
		if info["score"] != "*" {
			fmt.Println("Game Over:")
			fmt.Println(info["score"])
			break AiLoop
		}
	}

}

func aiVsHuman(game ghess.Board, reader *bufio.Reader) {
HumLoop:
	for {
		now := time.Now()
		state, err := ghess.MiniMax(0, 3, ghess.GetState(&game))
		if err != nil {
			fmt.Println(err)
			break HumLoop
		}
		fmt.Printf("\nThis took me: %s\n", time.Since(now))
		game.Move(state.Init[0], state.Init[1])
		info := game.Stats()
		fmt.Println(getPanel(info))
		fmt.Println(game.String())
		if info["score"] != "*" {
			fmt.Println("Game Over:")
			fmt.Println(info["score"])
			break HumLoop
		}
	InputLoop:
		for {
			fmt.Print("Your move: ")
			input, _ := reader.ReadString('\n')
			e := game.ParseMove(input)
			if e != nil {
				fmt.Println(e)
			} else {
				break InputLoop
			}
		}
		fmt.Println(getPanel(info))
		if info["turn"] == "b" {
			fmt.Println(game.StringBlack())
		} else {
			fmt.Println(game.StringWhite())
		}

		if info["score"] != "*" {
			fmt.Println("Game Over:")
			fmt.Println(info["score"])
			break HumLoop
		}
	}

}

func predictAiMove(game ghess.Board) {
	done := make(chan bool)
	go func() {
		state, err := ghess.MiniMax(0, 3, ghess.GetState(&game))
		fmt.Println(state)
		if err != nil {
			fmt.Println(err)
		}
		done <- true
	}()
	now := time.Now()
ThinkLoop:
	for {
		select {
		case <-done:
			fmt.Printf("\nThis took me: %s\n", time.Since(now))
			break ThinkLoop
		default:
			think(true)
		}

	}

}

func makeAiMove(game *ghess.Board) {
	now := time.Now()
	state, err := ghess.MiniMaxPruning(0, 3, ghess.GetState(game))
	if err != nil {
		fmt.Println(err)
	}
	game.Move(state.Init[0], state.Init[1])
	fmt.Println(game.String())
	fmt.Println(state)
	fmt.Printf("\nThis took me: %s\n", time.Since(now))
}

func getPanel(m map[string]string) string {
	return "|Move:  " + m["move"] + "     Turn: " + m["turn"] +
		"\n|Check: " + m["check"] + " Castle: " + m["castling"] +
		"\n|Mate:  " + m["checkmate"] + " Score: " + m["score"] + "\n"
}

func think(long bool) {
	var t int // time.Duration is no int
	if long {
		t = 500
	} else {
		t = 200
	}
	time.Sleep(time.Duration(t) * time.Millisecond)
	fmt.Print(".")
	time.Sleep(time.Duration(t) * time.Millisecond)
	fmt.Print(".")
	time.Sleep(time.Duration(t) * time.Millisecond)
	fmt.Print(".")
	time.Sleep(time.Duration(t) * time.Millisecond)
	fmt.Print(".")
}
