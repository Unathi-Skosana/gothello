package gomcts

type OthelloBoardGameAction struct {
  move   int8
  value  int8
}

type OthelloGameState struct {
	nextToMove   int8
	board        []int8
	emptySquares uint16
	ended        bool
	result       GameResult
}


// IsGameEnded - TicTacToeGameState implementation of IsGameEnded method of GameState interface
func (s OthelloGameState) IsGameEnded() bool {
	_, ended := s.EvaluateGame()
	return ended
}

// EvaluateGame - TicTacToeGameState implementation of EvaluateGame method of GameState interface
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
    if (s.board[i] == BLACK) {
      blackSum +=  1
    } else if (s.board[i] == WHITE) {
      whiteSum +=  1
    }
  }

  cur := s.GetLegalActions()[0].(OthelloBoardGameAction).move
  s.nextToMove *= -1
  next := s.GetLegalActions()[0].(OthelloBoardGameAction).move
  s.nextToMove *= -1


  if s.emptySquares == 0 || ( cur == -1 && next == -1) {
    if blackSum > whiteSum {
      return GameResult(1), true
    } else if blackSum < whiteSum {
      return GameResult(-1), true
    }
    return GameResult(0), true
  }

  return GameResult(0), false
}




// ApplyTo - TicTacToeBoardGameAction implementation of ApplyTo method of Action interface
func (a OthelloBoardGameAction) ApplyTo(s GameState) GameState {
  OthelloGameState := s.(OthelloGameState)
  OthelloGameState.board = copyOthelloBoard(OthelloGameState.board)

  if (a.move == -1) {
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

// GetLegalActions - TicTacToeGameState implementation of GetLegalActions method of GameState interface
func (s OthelloGameState) GetLegalActions() []Action {
  cnt := 0
  actions := make([]Action, 0, 0)
  for i := 11; i <= 88; i++ {
    if legalPlayer(int8(i), s) == 1 {
      cnt++
      actions = append(actions, OthelloBoardGameAction{move : int8(i), value: s.nextToMove})
    }
  }

  if cnt == 0 {
      actions = append(actions, OthelloBoardGameAction{move : -1, value: s.nextToMove})
  }

  return actions
}

// NextToMove - TicTacToeGameState implementation of NextToMove method of GameState interface
func (s OthelloGameState) NextToMove() int8 {
  return s.nextToMove
}
