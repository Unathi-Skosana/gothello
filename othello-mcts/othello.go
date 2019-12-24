package gomcts

import "fmt"

// OthelloBoardGameAction - action on a othello board
type OthelloBoardGameAction struct {
	move  int8
	value int8
}

// OthelloGameState - Othello game state
type OthelloGameState struct {
	nextToMove   int8
	board        []int8
	emptySquares uint16
	ended        bool
	result       GameResult
}

// PrintBoard - prints out game board to console
func (s OthelloGameState) PrintBoard() {
	fmt.Printf("   1 2 3 4 5 6 7 8 [%s=%d %s=%d]\n", nameof(BLACK), count(BLACK, s.board), nameof(WHITE), count(WHITE, s.board))
	for row := 1; row <= 8; row++ {
		fmt.Printf("%d  ", row)
		for col := 1; col <= 8; col++ {
			fmt.Printf("%s ", nameof(s.board[col+(10*row)]))
		}
		fmt.Printf("\n")
	}
}

// CreateOthelloInitialGameState - initializes a othello game state
func CreateOthelloInitialGameState() OthelloGameState {
	board := initializeOthelloBoard()
	state := OthelloGameState{nextToMove: BLACK, board: board, emptySquares: uint16(REALBOARDSIZE) - 4}
	return state
}

// IsGameEnded - OthelloGameState implementation of IsGameEnded method of GameState interface
func (s OthelloGameState) IsGameEnded() bool {
	_, ended := s.EvaluateGame()
	return ended
}

// EvaluateGame - OthelloGameState implementation of EvaluateGame method of GameState interface
func (s OthelloGameState) EvaluateGame() (result GameResult, ended bool) {

	defer func() {
		s.result = result
		s.ended = ended
	}()

	if s.ended {
		return s.result, s.ended
	}

	whiteSum := 0
	blackSum := 0

	for i := 1; i <= 88; i++ {
		if s.board[i] == BLACK {
			blackSum++
		} else if s.board[i] == WHITE {
			whiteSum++
		}
	}

	cur := s.GetLegalActions()[0].(OthelloBoardGameAction).move
	s.nextToMove *= -1
	next := s.GetLegalActions()[0].(OthelloBoardGameAction).move
	s.nextToMove *= -1

	if s.emptySquares == 0 || (cur == -1 && next == -1) {
		if blackSum > whiteSum {
			return GameResult(1), true
		} else if blackSum < whiteSum {
			return GameResult(-1), true
		}
		return GameResult(0), true
	}

	return GameResult(0), false
}

func (a OthelloBoardGameAction) GetMove() int8 {
	return a.move
}

// ApplyTo - OthelloBoardGameAction implementation of ApplyTo method of Action interface
func (a OthelloBoardGameAction) ApplyTo(s GameState) GameState {
	OthelloGameState := s.(OthelloGameState)
	OthelloGameState.board = copyOthelloBoard(OthelloGameState.board)

	if a.move == -1 {
		OthelloGameState.nextToMove *= -1
		return OthelloGameState
	}

	if OthelloGameState.nextToMove != a.value {
		panic("*hands slapped*,  not your turn")
	}

	makeMove(a.move, OthelloGameState)
	OthelloGameState.nextToMove *= -1
	OthelloGameState.emptySquares--

	return OthelloGameState
}

// GetLegalActions - OthelloGameState implementation of GetLegalActions method of GameState interface
func (s OthelloGameState) GetLegalActions() []Action {
	cnt := 0
	actions := make([]Action, 0, 0)
	for i := 11; i <= 88; i++ {
		if legalPlayer(int8(i), s) == 1 {
			cnt++
			actions = append(actions, OthelloBoardGameAction{move: int8(i), value: s.nextToMove})
		}
	}

	if cnt == 0 {
		actions = append(actions, OthelloBoardGameAction{move: -1, value: s.nextToMove})
	}

	return actions
}

// NextToMove - OthelloGameState implementation of NextToMove method of GameState interface
func (s OthelloGameState) NextToMove() int8 {
	return s.nextToMove
}

// ToString - return game board as string
func (s OthelloGameState) ToString() string {
	board := s.board
	boardStr := ""

	for row := 1; row <= 8; row++ {
		for col := 1; col <= 8; col++ {
			boardStr += nameof(board[col+(10*row)])
		}
	}

	return boardStr
}
