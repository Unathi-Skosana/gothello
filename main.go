package main

import (
    "othello"
)

func main() {
  var state gomcts.GameState = gomcts.CreateOthelloInitialGameState(1)
  gomcts.PrintBoard(state.(gomcts.OthelloGameState))
  for !state.IsGameEnded() {
    chosenAction:= gomcts.MonteCarloTreeSearch(state, gomcts.OthelloHeuristicRolloutPolicy, 100)
    //fmt.Println(chosenAction)
    state = chosenAction.ApplyTo(state)
    gomcts.PrintBoard(state.(gomcts.OthelloGameState))
  }
}
