package othello

import (
	"crypto/rand"
	"math/big"

	"github.com/unathi-skosana/gothello/pkg/gomcts"
)

// OthelloRandomRolloutPolicy - Randomly select next move
func OthelloRandomRolloutPolicy(state gomcts.GameState) gomcts.Action {
	actions := state.GetLegalActions()
	numberOfActions := len(actions)

	if numberOfActions == 1 {
		return actions[0]
	}

	actionIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(numberOfActions)))
	return actions[actionIndex.Int64()]

}

// OthelloMediumRolloutPolicy - Evaluate moves with evaluation function and
// select one with max evaluation score with equally weighted heuristics
func OthelloMediumRolloutPolicy(state gomcts.GameState) gomcts.Action {
	actions := state.GetLegalActions()
	scores := make([]float64, 0)
	dummyGameState := state.(OthelloGameState)
	numberOfActions := len(actions)

	parityWeight := 25.00
	mobilityWeight := 25.00
	cornersWeight := 25.00
	frontiersWeight := 25.00

	if numberOfActions == 1 {
		return actions[0]
	}

	for i := 0; i < numberOfActions; i++ {
		cur := actions[i].ApplyTo(dummyGameState.Clone())
		scores = append(scores, evaluate(cur.(OthelloGameState), parityWeight, mobilityWeight, cornersWeight, frontiersWeight))
	}

	maxIndex := 0
	maxValue := scores[0]

	for i := 1; i < numberOfActions; i++ {
		if scores[i] > maxValue {
			maxValue = scores[i]
			maxIndex = i
		}
	}

	return actions[maxIndex]
}

// OthelloHardRolloutPolicy - Evaluate moves with evaluation function and select
// one with max evaluation score with heuristic weights I've found to work quite
// work.
func OthelloHardRolloutPolicy(state gomcts.GameState) gomcts.Action {
	actions := state.GetLegalActions()
	scores := make([]float64, 0)
	dummyGameState := state.(OthelloGameState)
	numberOfActions := len(actions)

	parityWeight := 21.45
	mobilityWeight := 3.37
	cornersWeight := 70.00
	frontiersWeight := 5.38

	if numberOfActions == 1 {
		return actions[0]
	}

	for i := 0; i < numberOfActions; i++ {
		cur := actions[i].ApplyTo(dummyGameState.Clone())
		scores = append(scores, evaluate(cur.(OthelloGameState), parityWeight, mobilityWeight, cornersWeight, frontiersWeight))
	}

	maxIndex := 0
	maxValue := scores[0]

	for i := 1; i < numberOfActions; i++ {
		if scores[i] > maxValue {
			maxValue = scores[i]
			maxIndex = i
		}
	}

	return actions[maxIndex]
}
