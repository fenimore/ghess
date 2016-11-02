package ghess

import (
	"errors"
	"fmt"
	"strings"
)

// Principal Variation Search

var pvHash map[[120]byte]int = make(map[[120]byte]int)

/*
MiniMax implementation ###########################################
*/

// State struct holds a board position,
// the move that got there, and the evaluation.
// Init is the move which began a certain branch of the tree.
type State struct {
	board *Board // the Board object
	eval  int    // score
	Init  [2]int // the moves which got to that position at root
	isMax bool   // is White Player
	alpha int
	beta  int
}

// String returns some basic info of a State.
func (s State) String() string {
	return fmt.Sprintf("\nScore: %d\nFrom Move: %d, %d", s.eval, s.Init[0], s.Init[1])
}

// States are a slice of State structs.
type States []State

// Sort functionality depricated.
func (s States) Len() int           { return len(s) }
func (s States) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s States) Less(i, j int) bool { return s[i].eval < s[j].eval }

// GetState turns a Board into a copy and it's state.
// The Init value is nil.
func GetState(b *Board) State {
	c := *b // dereference the pointer
	s := State{board: &c, eval: c.Evaluate()}
	return s
}

// TryState takes in a *Board and valid move and returns
// a State struct.
func TryState(b *Board, o, d int) (State, error) {
	state := State{}
	possible := CopyBoard(b)
	err := possible.Move(o, d)
	if err != nil {
		return state, err
	}
	state.board = possible
	state.eval = possible.Evaluate()
	return state, nil
}

// GetPossibleStates returns a slice of State structs
// Each with a score and the move that got there.
func GetPossibleStates(state State) (States, error) {
	states := make(States, 0)
	origs, dests := state.board.SearchValid()
	for i := 0; i < len(origs); i++ {
		s, err := TryState(state.board, origs[i], dests[i])
		if err != nil {
			return states, err
		}
		if state.Init[0] == 0 {
			s.Init[0], s.Init[1] = origs[i], dests[i]
		} else {
			s.Init[0], s.Init[1] = state.Init[0], state.Init[1]
		}
		s.isMax = state.isMax // Basically is White
		states = append(states, s)
	}
	return states, nil
}

// GetPossibleStates returns a slice of State structs
// Each with a score and the move that got there.
func GetPossibleOrderedStates(state State) (States, error) {
	pv, ok := pvHash[state.board.board]
	if ok {
		fmt.Println(pv)
	}
	states := make(States, 0)
	origs, dests := state.board.SearchValid() //.SearchValidOrdered()
	for i := 0; i < len(origs); i++ {
		s, err := TryState(state.board, origs[i], dests[i])
		if err != nil {
			return states, err
		}
		if state.Init[0] == 0 {
			s.Init[0], s.Init[1] = origs[i], dests[i]
		} else {
			s.Init[0], s.Init[1] = state.Init[0], state.Init[1]
		}
		s.isMax = state.isMax // Basically is whitePlayer or !whitePlayer
		//pvHash[s.board.board] = s.eval
		states = append(states, s)
	}
	return states, nil
}

// DictionaryAttack looks up common openings
// for less stupid opening moves.
func DictionaryAttack(s State) (State, error) {
	position := s.board.Position()
	// I don't need castling empassant or move number
	posits := strings.Split(position, " ")
	key := posits[0] + " " + posits[1]
	// Check if opening exists
	if val, ok := dict[key]; ok {
		state := State{Init: val}
		return state, nil
	}

	return s, errors.New("No Dictionary Attack Found")
}

// The Max and Mini functions are O(n)

// Max returns the state from States with
// highest evaluation.
func Max(states States) State {
	var maxIdx int
	var maxVal int = -1000
	for idx, state := range states {
		if state.eval > maxVal {
			maxVal = state.eval
			maxIdx = idx
		}
	}
	return states[maxIdx]
}

// Min returns the state from States with
// Lowest evaluation.
func Min(states States) State {
	var minIdx int
	var minVal int = 1000
	for idx, state := range states {
		if state.eval < minVal {
			minVal = state.eval
			minIdx = idx
		}
	}
	return states[minIdx]
}

// MiniMax Recursive, pass in state, search depth and terminal depth.
// and depth is always 0 when passed in initially.
// This is like a DFS algorithm which tries to Minimize maximun loss.
// TODO: write tests somehow.
// Pass bback error LOL
func MiniMaxPruning(depth, terminal int, s State) (State, error) {
	if depth == 0 {
		s.alpha = -1000000000
		s.beta = 1000000000
		// set the Min or Max
		if s.board.toMove == "w" {
			s.isMax = true
		} else {
			s.isMax = false
		}
		openState, err := DictionaryAttack(s)
		if err == nil {
			return openState, nil
		}
	}
	if depth == terminal {
		return s, nil
	}

	even := (depth % 2) == 0
	var maxNode bool
	if even {
		// If White Player Return Maximum
		if s.isMax {
			maxNode = true
		} else {
			maxNode = false
		}
	} else { // Otherwise Return Minimum... Yup that's the idea.
		if s.isMax {
			maxNode = false
		} else {
			maxNode = true
		}
	}

	states, err := GetPossibleStates(s)
	if err != nil {
		return s, err
	}

	// Recursive Call
	var bestState State
	var bestStates States
	for _, state := range states {
		state.alpha = s.alpha
		state.beta = s.beta
		bestState, err = MiniMaxPruning(depth+1, terminal, state)
		if err != nil {
			return bestState, err
		}
		// The trick is to update the root (for this branch)
		// beta or alpha, which then will cut off further iterating in
		// THIS VERY for loop.
		if maxNode {
			if bestState.eval > s.beta {
				//fmt.Println("Bingo Alpha", bestState.eval)
				return bestState, nil
			} else {
				bestState.beta = s.beta
				s.alpha = max(s.alpha, bestState.eval)
			}
		}
		if !maxNode {
			if bestState.eval < s.alpha {
				//fmt.Println("Bingo Beta", bestState.eval)
				return bestState, nil
			} else {
				bestState.alpha = s.alpha
				s.beta = min(s.beta, bestState.eval)

			}
		}

		// If the player is Max, I want to compare against beta
		// otherwise against alpha.

		// If we are considering Max,
		//and state's value >= beta, then return NOW
		// otherwise, set alpha = Max(alpha, state's value)

		// If we are considering Min,
		// and state's value <= alpha, then return NOW
		// otherwise, set beta = Min(beta, state's value)

		bestStates = append(bestStates, bestState)
	}
	if len(bestStates) < 1 {
		return s, nil
	}

	if maxNode {
		return Max(bestStates), nil
	} else {
		return Min(bestStates), nil
	}
}

// MinimaxOrdered prunes an ordered list of states
func MiniMaxOrdered(depth, terminal int, s State) (State, error) {
	if depth == 0 {
		s.alpha = -1000000000
		s.beta = 1000000000
		// set the Min or Max
		if s.board.toMove == "w" {
			s.isMax = true
		} else {
			s.isMax = false
		}
		//fmt.Println("SHHH, I'm thinking")
		// DICT attack
		openState, err := DictionaryAttack(s)
		if err == nil {
			return openState, nil
		}
	}
	if depth == terminal { // that is, 2 ply
		return s, nil
	}

	even := (depth % 2) == 0
	var maxNode bool
	if even {
		// If White Player Return Maximum
		if s.isMax {
			maxNode = true
			//return Max(bestStates), nil
		} else {
			maxNode = false
			//return Min(bestStates), nil
		}
	} else { // Otherwise Return Minimum... Yup that's the idea.
		if s.isMax {
			maxNode = false
			//return Min(bestStates), nil
		} else {
			maxNode = true
			//return Max(bestStates), nil
		}
	}

	states, err := GetPossibleOrderedStates(s)
	if err != nil {
		return s, err
	}

	// Recursive Call
	var bestState State
	var bestStates States
	for _, state := range states {
		pv, ok := pvHash[s.board.board]
		if ok {
			fmt.Println("This Should exist", pv)
		}
		state.alpha = s.alpha
		state.beta = s.beta

		bestState, err = MiniMaxPruning(depth+1, terminal, state)
		if err != nil {
			return bestState, err
		}
		// The trick is to update the root (for this branch)
		// beta or alpha, which then will cut off further iterating in
		// THIS VERY for loop.
		// s is root
		if maxNode {
			if bestState.eval > s.beta {
				//fmt.Println("Bingo Alpha", bestState.eval)
				return bestState, nil
			} else {
				bestState.beta = s.beta
				pvHash[bestState.board.board] = bestState.eval
				s.alpha = max(s.alpha, bestState.eval)
			}
		}
		if !maxNode {
			if bestState.eval < s.alpha {
				//fmt.Println("Bingo Beta", bestState.eval)
				//pvHash[bestState.board.board] = bestState.eval
				return bestState, nil
			} else {
				pvHash[bestState.board.board] = bestState.eval
				bestState.alpha = s.alpha
				s.beta = min(s.beta, bestState.eval)

			}
		}

		// If the player is Max, I want to compare against beta
		// otherwise against alpha.

		// If we are considering Max,
		//and state's value >= beta, then return NOW
		// otherwise, set alpha = Max(alpha, state's value)

		// If we are considering Min,
		// and state's value <= alpha, then return NOW
		// otherwise, set beta = Min(beta, state's value)

		bestStates = append(bestStates, bestState)
	}
	if len(bestStates) < 1 {
		return s, nil
	}

	if maxNode {
		return Max(bestStates), nil
	} else {
		return Min(bestStates), nil
	}
}

// small min, doesn't take state,
// but it takes numbers
func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// DEPRECATED

// Depricated
func MiniMax(depth, terminal int, s State) (State, error) {
	if depth == 0 {
		// set the Min or Max
		if s.board.toMove == "w" {
			s.isMax = true
		} else {
			s.isMax = false
		}
		//fmt.Println("SHHH, I'm thinking")
		// DICT attack
		openState, err := DictionaryAttack(s)
		if err == nil {
			return openState, nil
		}
	}
	if depth == terminal { // that is, 2 ply
		//fmt.Println("Depth ", depth, s)
		return s, nil
	}
	// Determine which node this is
	// TODO: Why is this so complicated?
	// Because when minimax from perspective black,
	// things are totally different
	even := (depth % 2) == 0
	var maxNode bool
	if even {
		// If White Player Return Maximum
		if s.isMax {
			maxNode = true
			//return Max(bestStates), nil
		} else {
			maxNode = false
			//return Min(bestStates), nil
		}
	} else { // Otherwise Return Minimum... Yup that's the idea.
		if s.isMax {
			maxNode = false
			//return Min(bestStates), nil
		} else {
			maxNode = true
			//return Max(bestStates), nil
		}
	}

	states, err := GetPossibleStates(s)
	if err != nil {
		return s, err
	}

	// Recursive Call
	var bestState State
	var bestStates States
	for _, state := range states {
		bestState, err = MiniMax(depth+1, terminal, state)
		if err != nil {
			return bestState, err
		}
		bestStates = append(bestStates, bestState)
	}
	if len(bestStates) < 1 {
		return s, nil
	}

	if maxNode {
		return Max(bestStates), nil
	} else {
		return Min(bestStates), nil
	}
}
