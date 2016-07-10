package main

import "github.com/polypmer/ghess"
import "fmt"
import "bufio"
import "regexp"
import "os"
import "strconv"
import "strings"


func main() {
	game := ghess.NewBoard()
	PlayGame(game)
}

func PlayGame(game ghess.Board) { // TODO Rotate Board
	var info map[string]string
	var turn string
	welcome := `
********
go-chess
    Enter /help for more options

    /~ |_ _  _ _
    \_|||(/__\_\

`
	manuel := `Help:
    Prefix commands with / - slash

Commands:
	quit - exit game
	new - new game
        print - print game
        panel - print game info
	coordinates - print board coordinates
	pgn - print PGN history
	fen - print FEN position
	set-headers - set PGN headers
	headers - print game info
Tests:
	test-castle
        test-pgn - load a pgn game
`
	reader := bufio.NewReader(os.Stdin)
	// welcome message
	fmt.Println(welcome)
	fmt.Print(game.String())
	info = game.Stats()
Loop:
	for {
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
			switch {
			case input == "/help":
				fmt.Print("\n", manuel)
			case input == "/quit":
				break Loop 
			case input == "/new":
				game = ghess.NewBoard()//not a Board method
				fmt.Print(game.String())
			case input == "/print":
				fmt.Print(game.String())
			case input == "/panel":
				info = game.Stats()
				fmt.Print(getPanel(info))
			case input == "/coordinates":
				fmt.Println("Coordinates:")
				game.Coordinates()
			case input == "/pgn":
				fmt.Println("PGN history:")
				fmt.Println(game.PgnString())
			case input == "/load-pgn":
				var err error
				fmt.Print("Enter PGN history: ")
				history, _ := reader.ReadString('\n')
				game, err = game.LoadPgn(history)
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
				fmt.Println(info["headers"])
			case input == "/test-pgn":
				hist := `1. Nf3 Nc6 2. d4 d5 3. c4 e6 4. e3 Nf6 5. Nc3 Be7 6. a3 O-O 7. b4 a6 8. Be2 Re8 9. O-O Bf8 10. c5 g6 11. b5 axb5 12. Bxb5 Bd7 13. h3 Na5 14. Bd3 Nc6 15. Rb1 Qc8 16. Nb5 e5 17. Be2 e4 18. Ne1 h6 19. Nc2 g5 20. f3 exf3 21. Bxf3 g4 22. hxg4 Bxg4 23. Nxc7 Qxc7 24. Bxg4 Nxg4 25. Qxg4+ Bg7 26. Nb4 Nxb4 27. Rxb4 Ra6 28. Rf5 Re4 29. Qh5 Rg6 30. Qh3 Qc8 31. Qf3 Qd7 32. Rb2 Bxd4 33. exd4 Re1+ 34. Kh2 Rxc1 35. Qxd5 Qe7 36. g3 Qc7 37. Rf4 b6 38. a4 Rg5 39. cxb6 Rxd5 40. bxc7 Rxc7 41. Rb5 Rc2+`
				var err error
				game, err = game.LoadPgn(hist)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(game.String())
				ch,_:= strconv.ParseBool(info["check"])
				if ch {
					fmt.Println("****Check!****")
				}
			case input == "/test-check":
				hist := `1. e4 e5 2. Qf3 Qg5 3. Qxf7`
				var err error
				game, err = game.LoadPgn(hist)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(game.String())
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
		fmt.Print(game.String())
		ch,_:= strconv.ParseBool(info["check"])
		if ch {
			fmt.Println("****Check!****")
		}
	}
	fmt.Println("\nGood Game.")
}

func getPanel(m map[string]string) string {
	return "|Debug Pane:\n|Move: "+m["move"]+" Turn: "+m["turn"]+"\n|Check: "+m["check"]+" Castle: "+m["castling"]+"\n"
}
