package othello

import (
	"testing"

	"github.com/unathi-skosana/gothello/pkg/gomcts"
)

func TesOthelloGameStateInitialization(t *testing.T) {
	state := New(1)
	board := state.GetBoard()

	if state.nextToMove != 1 {
		t.Errorf("state.nextToMove should be 1, but it is %v", state.nextToMove)
	}

	if len(board) != 100 {
		t.Errorf("len(state.board) should be 100, but it is %v", len(state.board))
	}

	if board[54] == BLUE && board[45] == BLUE && board[44] == RED && board[55] == RED {
		t.Errorf("pieces at 54, 45, 44 and 5  should be 1, 1, 2 and 1,  but they are %v, %v, %v and  %v", board[54], board[45], board[44], board[55])
	}
}

func TestMoveProducesOthelloGameStateCorrectly(t *testing.T) {
	state := New(1)
	action := OthelloBoardGameAction{move: 34, value: 1}
	nextState := action.ApplyTo(state).(OthelloGameState)

	if nextState.nextToMove != 2 {
		t.Errorf("state.nextToMove should be 2, but it is %v", state.nextToMove)
	}

	if &(nextState.board[34]) == &(state.board[34]) {
		t.Errorf("state.board[0][0] and nextState.board[0][0] refer to the same memory location - but should not")
	}

	if nextState.board[34] != 1 {
		t.Errorf("nextState.board[1][1] should be 1 but is %v", nextState.board[34])
	}

	if state.board[34] != 0 {
		t.Errorf("state.board[34] should remain 0 but is %v", state.board[34])
	}

}

func TestEmptyOthelloGameStateEvaluation(t *testing.T) {
	state := New(1)
	_, gameEnded := state.EvaluateGame()
	if gameEnded {
		t.Errorf("Game state is evaluated as ended but should not")
	}
}

func TestNumberOfLegalActionsOfOthelloGameState(t *testing.T) {
	state := New(1)
	actions := state.GetLegalActions()
	if len(actions) != 4 {
		t.Errorf("There should be 4 actions to perform but is %v", len(actions))
	}

}

func TestOthellolGameStateZeroIfGameEnded(t *testing.T) {
	state := New(1)
	board := make([]int, BOARD_SIZE)
	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			board[i] = RED
		}
	}

	state.board = board

	actions := state.GetLegalActions()
	if len(actions) > 0 {
		t.Errorf("Game is ended but state still has actions")
	}
}

func TestNotYourTurnPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic but should")
		}
	}()

	state := New(1)
	action := OthelloBoardGameAction{move: 34, value: 2}
	action.ApplyTo(state)
}

func TestOutOfBoardMovePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic but should")
		}
	}()

	state := New(1)
	action := OthelloBoardGameAction{move: 99, value: 1}
	action.ApplyTo(state)
}

func TestAlreadyOccupiedSquareMovePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic but should")
		}
	}()

	state := New(1)
	action := OthelloBoardGameAction{move: 44, value: 1}
	action.ApplyTo(state)
}

func TestGameEvaluationShouldBeNotEnded(t *testing.T) {
	state := New(1)

	_, ended := state.EvaluateGame()

	if ended {
		t.Errorf("Game should be not ended but is")
	}

	if state.ended {
		t.Errorf("Game should be not ended but is")
	}

}

func TestGameEvaluationShouldBeDraw(t *testing.T) {
	state := New(1)
	board := make([]int, BOARD_SIZE)
	j := 0

	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			if j == 0 {
				board[i] = RED
			} else {
				board[i] = BLUE
			}
			j = (j + 1) % 2
		}
	}

	state.board = board

	result, ended := state.EvaluateGame()
	if result != gomcts.GameResult(0) {
		t.Errorf("Result should be a draw but is %v", result)
	}

	if !ended {
		t.Errorf("Game should be ended but is not")
	}

	if state.ended {
		t.Errorf("Game should be ended but is not")
	}
}

func TestGameEvaluationShouldResultBlueWinning(t *testing.T) {
	state := New(1)
	board := make([]int, BOARD_SIZE)
	j := 0

	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			if j == 0 {
				board[i] = RED
			} else {
				board[i] = BLUE
			}
			j = (j + 1) % 3
		}
	}

	state.board = board

	result, ended := state.EvaluateGame()
	if result != gomcts.GameResult(BLUE) {
		t.Errorf("Result should be a 1 but is %v", result)
	}

	if !ended {
		t.Errorf("Game should be ended but is not")
	}

}

func TestGameEvaluationShouldResultRedWinning(t *testing.T) {
	state := New(1)
	board := make([]int, BOARD_SIZE)
	j := 0
	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			if j == 0 {
				board[i] = BLUE
			} else {
				board[i] = RED
			}
			j = (j + 1) % 3
		}
	}

	state.board = board

	result, ended := state.EvaluateGame()
	if result != gomcts.GameResult(RED) {
		t.Errorf("Result should be a 2 but is %v", result)
	}

	if !ended {
		t.Errorf("Game should be ended but is not")
	}

}
