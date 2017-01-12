package ghess

import (
	"errors"
	"fmt"
	"strings"
)

// Principal Variation Search

var pvHash map[[120]byte]int = make(map[[120]byte]int)
var pvMap map[Board][2]int = make(map[Board][2]int)

/*
MiniMax implementation ###########################################
*/

// State struct holds a board position,
// the move that got there, and the evaluation.
// Init is the move which began a certain branch of the tree.
type State struct {
	board  *Board // the Board object
	eval   int    // score
	Init   [2]int // the moves which got to that position at root
	isMax  bool   // is White Player
	alpha  int
	beta   int
	parent *State
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
		//fmt.Println(err, o, d)
		//fmt.Println(b.StringWhite())
		//fmt.Println(b.toMove)
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
			s.Init[0], s.Init[1] =
				state.Init[0], state.Init[1]
		}
		s.isMax = state.isMax // Basically is White
		// Add parent state?
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

// MiniMax with Alpha Beta Pruning:
// Return the position where, assuming your opponent picks its
// *best* moves, they have the minimum advantage
// Minimizing Maximum Loss
//
// Param:
//     state, current depth and terminal depth.
//     Depth is always 0 when passed in initially.
//
// Minimax:
//     This is like a Depth First Search algorithm.
//     Speed increase from Pruning largely depends on Move Ordering
func MiniMaxPruning(depth, terminal int, s State) (State, error) {
	if depth == 0 {
		// At first depth set Alpha and Beta values
		s.alpha = -1000000000
		s.beta = 1000000000

		// Search for Max or Min Player?
		if s.board.toMove == "w" {
			s.isMax = true
		} else {
			s.isMax = false
		}

		// At first depth check for Opening in Dictionary
		openState, err := DictionaryAttack(s)
		if err == nil {
			return openState, nil
		}
	}

	if depth == terminal {
		// The final state will pass up the
		// call stack:
		// the Score of terminal position
		// the Original move to get there
		return s, nil
	}

	// Determine if Max or Min node
	// by height of tree
	even := (depth % 2) == 0
	maxNode := even == s.isMax

	states, err := GetPossibleStates(s)
	if err != nil {
		return s, err
	}

	// Recursively call MiniMax on all Possible States
	var bestState State
	var bestStates States
	for _, state := range states {
		state.alpha = s.alpha
		state.beta = s.beta
		// Increment Depth when calling MiniMax
		bestState, err = MiniMaxPruning(depth+1, terminal, state)
		if err != nil {
			return bestState, err
		}

		/* Alpha Beta Pruning
		* The trick is to update the root node's
		* (for this branch) beta or alpha
		* This will prune further branches from root

		* If we are considering a Max Node,
		* and state's score/evaluation >= beta, then PRUNE/return
		* otherwise, set alpha = Max(alpha, state's value)

		* If we are considering a Min Node,
		* and state's score/evaluation <= alpha, then PRUNE/return
		* otherwise, set beta = Min(beta, state's value)

		* Pruning effectively breaks out of this loop
		* which call MiniMax on all Possible States
		* in a certain branch of the tree
		 */
		if maxNode {
			if bestState.eval > s.beta {
				return bestState, nil
			} else {
				bestState.beta = s.beta
				s.alpha = max(s.alpha, bestState.eval)
			}
		}
		if !maxNode {
			if bestState.eval < s.alpha {
				return bestState, nil
			} else {
				bestState.alpha = s.alpha
				s.beta = min(s.beta, bestState.eval)

			}
		}

		// if there is no Pruning, add the returned
		// bestStates from MiniMax calls to slice of bestStates
		bestStates = append(bestStates, bestState)
	}

	// If say there is stalemate/checkmate no best-states
	if len(bestStates) < 1 {
		return s, nil
	}

	if maxNode { // if height == Max nodes
		return Max(bestStates), nil
	} else { // if height == Min nodes
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
// MiniMax is Deprecated in favor of
// AlphaBeta Pruning Minimax
// Use for Testing
func MiniMax(depth, terminal int, s State) (State, error) {
	if depth == 0 {
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
