package ghess

import "fmt"

/*
MiniMax implementation ###########################################
*/

// State struct holds a board position,
// the move that got there, and the evaluation.
// Init is the move which began a certain branch of the tree.
type State struct {
	board *Board
	eval  int
	Init  [2]int
	isMax bool // find for maximun player
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
	//	if depth == 1 {
	//		fmt.Println(len(states), s.eval)
	//	}
	//fmt.Println(len(states), s.eval, depth

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
	//	if depth == 1 {
	//		fmt.Println(len(states), s.eval)
	//	}
	//fmt.Println(len(states), s.eval, depth

	for _, state := range states {
		bestState, err = MiniMax(depth+1, terminal, state)
		if err != nil {
			return bestState, err
		}
		bestStates = append(bestStates, bestState)

		// beta == +10000

		// If the player is Max, I want to compare against beta
		// otherwise against alpha.

		// If we are considering Max,
		//and state's value >= beta, then return NOW
		// otherwise, set alpha = Max(alpha, state's value)

		// If we are considering Min,
		// and state's value <= alpha, then return NOW
		// otherwise, set beta = Min(beta, state's value)
		if maxNode {
			if bestState.eval < s.alpha {
				fmt.Println("Alpha")
				return bestState, nil
			} else {
				if bestState.eval < s.beta {
					bestState.beta = bestState.eval
				}
			}
		}
		if !maxNode {
			if bestState.eval > s.beta {
				fmt.Println("BETA")
				return bestState, nil
			} else {
				if bestState.eval > s.alpha {
					bestState.alpha = bestState.eval
				}
			}
		}

		/*		if !maxNode {
					if bestState.eval <= s.alpha {
						fmt.Println("ALPHA")
						return bestState, nil
					} else {
						if bestState.eval < s.beta {
							bestState.beta = bestState.eval
						}
					}
				}
				if maxNode {
					if bestState.eval >= s.beta {
						fmt.Println("BETA")
						return bestState, nil
					} else {
						if bestState.eval > s.alpha {
							bestState.alpha = bestState.eval
						}
					}
				}*/

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
