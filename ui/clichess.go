package main

import "github.com/polypmer/ghess"
import "fmt"
import "bufio"
import "regexp"
import "os"
import "strconv"
import "strings"
import "time"

func main() {
	game := ghess.NewBoard()
	PlayGame(game)
}

// PlayGame is the command line user interface.
// It is used mostly for debugging.
func PlayGame(game ghess.Board) { // TODO Rotate Board
	var info map[string]string
	var turn string
	welcome := `
=================
         go-chess

    /~ |_ _  _ _
    \_|||(/__\_\


    Enter /help for more options

`
	manuel := `
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
    ai           - watch AI match
    vs           - play against AI as White

`
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
			// REFACTOR:
			// Have Switch (input)... simplify
			switch {
			case input == "/help":
				fmt.Print("\n", manuel)
			case input == "/quit" || input == "/exit":
				break Loop
			case input == "/new":
				game = ghess.NewBoard()
				fmt.Print(game.String())
			case input == "/print":
				fmt.Print(game.String())
			case input == "/panel":
				info = game.Stats()
				fmt.Print(getPanel(info))
			case input == "/coordinates":
				fmt.Println("Coordinates:")
				game.Coordinates()
			case input == "/score":
				checkMate, _ := strconv.ParseBool(info["checkmate"])
				score := info["score"]
				if checkMate { // TODO: or draw
					fmt.Println("Game over")
				}
				fmt.Println("Score: ", score)
			case input == "/pgn":
				fmt.Println("PGN history:")
				fmt.Println(game.PgnString())
			case input == "/load-pgn":
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
			case input == "/load-fen":
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
			case input == "/fen":
				fmt.Println("FEN position:")
				fmt.Println(game.Position())
			case input == "/set-headers":
				fmt.Print("Enter White Player: ")
				first, _ := reader.ReadString('\n')
				fmt.Print("Enter Black Player: ")
				second, _ := reader.ReadString('\n')
				game.SetHeaders(first, second)
			case input == "/headers":
				info = game.Stats()
				fmt.Println(info["headers"])
			case input == "/valid":
				origs, dests := game.SearchValid()
				fmt.Println(origs)
				fmt.Println(dests)
				fmt.Println("Total valid moves: ",
					len(origs))
			case input == "/eval":
				score := game.Evaluate()
				fmt.Println("Position: ", score)
			case input == "/minimax":
				done := make(chan bool)
				go func() {
					state := ghess.MiniMax(0, 4, ghess.GetState(&game))
					fmt.Println(state)
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
			case input == "/ai":
				for {
					fmt.Println(game.StringWhite())
					state := ghess.MiniMax(0, 3, ghess.GetState(&game))
					game.Move(state.Init[0], state.Init[1])
				}
			case input == "/rand":
				origs, dests := game.SearchValid()
				e := game.MoveRandom(origs, dests)
				if e != nil {
					fmt.Println(e)
				}
			case input == "/random-game":
				reader := bufio.NewReader(os.Stdin)
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
					info = game.Stats()
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
			case input == "/tension":
				fmt.Println("Tension Coordinates:")
				fmt.Println(game.Tension())
				fmt.Println("Total Tension: ", game.TensionSum())
			default:
				fmt.Println("Mysterious input")
			}
			continue
		}
		e := game.ParseMove(input)
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
	fmt.Println("\nGood Game.")
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
