package main

import (
	"fmt"

	"github.com/unathi-skosana/gothello/gomcts"
	"github.com/unathi-skosana/gothello/othello"
)

func main() {

	var i int = 1
	var chosenAction gomcts.Action
	var s gomcts.GameState = othello.New()

	othello.PrintBoard(s)

	for !s.IsGameEnded() {
		if i%2 == 1 {
			chosenAction = gomcts.MonteCarloTreeSearch(s, othello.OthelloMediumRolloutPolicy, 1500)
		} else {
			chosenAction = gomcts.MonteCarloTreeSearch(s, othello.OthelloHardRolloutPolicy, 1500)
		}
		s = chosenAction.ApplyTo(s)
		othello.PrintBoard(s)
		fmt.Println(chosenAction)
		i++
	}
}
