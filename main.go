package main

import "github.com/unathi-skosana/gothello/src/gomcts"

func main() {
	var s gomcts.GameState = gomcts.CreateOthelloInitialGameState()
	s.(gomcts.OthelloGameState).PrintBoard()
	for !s.IsGameEnded() {
		chosenAction := gomcts.MonteCarloTreeSearch(s, gomcts.OthelloHeuristicRolloutPolicy, 1000)
		s = chosenAction.ApplyTo(s)
		s.(gomcts.OthelloGameState).PrintBoard()
	}
}
