package main

import "gomcts"

func main() {
	var state gomcts.GameState = gomcts.CreateOthelloInitialGameState()
	gomcts.PrintBoard(state.(gomcts.OthelloGameState))
	for !state.IsGameEnded() {
		chosenAction := gomcts.MonteCarloTreeSearch(state, gomcts.OthelloHeuristicRolloutPolicy, 1000)
		//fmt.Println(chosenAction)
		state = chosenAction.ApplyTo(state)
		state.PrintBoard()
	}
}
