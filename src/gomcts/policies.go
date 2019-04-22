package gomcts

import (
	"crypto/rand"
	"math/big"
)

// RolloutPolicy - function signature determining the next action during Monte Carlo Tree Search rollout
type RolloutPolicy func(GameState) Action

// OthelloRandomRolloutPolicy - Randomly select next move
func OthelloRandomRolloutPolicy(state GameState) Action {
	actions := state.GetLegalActions()
	numberOfActions := int64(len(actions))

	if numberOfActions == 1 {
		return actions[0]
	}

	actionIndex, _ := rand.Int(rand.Reader, big.NewInt(numberOfActions))
	return actions[actionIndex.Int64()]

}

// OthelloHeuristicRolloutPolicy - Evaluate moves with evaluation function and select one with max evaluation score
func OthelloHeuristicRolloutPolicy(state GameState) Action {
	actions := state.GetLegalActions()
	scores := make([]float64, 0)
	dummyGameState := state.(OthelloGameState)
	numberOfActions := int64(len(actions))

	if numberOfActions == 1 {
		return actions[0]
	}

	var i int64

	for i = 0; i < numberOfActions; i++ {
		cur := actions[i].ApplyTo(copyOthelloGameState(dummyGameState))
		scores = append(scores, evalFunc(cur.(OthelloGameState)))
	}

	var maxIndex int64
	var maxValue = scores[0]

	for i = 1; i < numberOfActions; i++ {
		if scores[i] > maxValue {
			maxValue = scores[i]
			maxIndex = i
		}
	}

	return actions[maxIndex]
}
