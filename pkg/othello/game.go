package othello

import (
	"github.com/unathi-skosana/gothello/pkg/gomcts"
)

// Pieces
const (
	EMPTY = iota
	BLUE
	RED
)

// DIR
const DIR = 8

// sizes
const BOARD_SIZE = 100
const BOARD_WIDTH = 8
const PIECE_SLOTS = 64

// BLOCKS
const FIRST_BLOCK = 11
const LAST_BLOCK = 88

// DIRECTIONS
var ALLDIRECTIONS = []int{-11, -10, -9, -1, 1, 9, 10, 11}

// OthelloBoardGameAction - action on a othello board
type OthelloBoardGameAction struct {
	move  int
	value int
}

// OthelloGameState - Othello game state
type OthelloGameState struct {
	nextToMove int
	board      []int
	ended      bool
	result     gomcts.GameResult
}

/*
 * GameState interface methods
 */

// New - initializes a new  OthelloGameState object
func New(nextToMove int) OthelloGameState {
	board := make([]int, BOARD_SIZE)
	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			board[i] = EMPTY
		}
	}

	board[44] = RED
	board[45] = BLUE
	board[54] = BLUE
	board[55] = RED

	state := OthelloGameState{nextToMove: nextToMove, board: board}
	return state
}

// IsGameEnded - OthelloGameState implementation of IsGameEnded method of GameState interface
func (s OthelloGameState) IsGameEnded() bool {
	_, ended := s.EvaluateGame()
	return ended
}

// EvaluateGame - OthelloGameState implementation of EvaluateGame method of GameState interface
func (s OthelloGameState) EvaluateGame() (result gomcts.GameResult, ended bool) {

	defer func() {
		s.result = result
		s.ended = ended
	}()

	if s.ended {
		return s.result, s.ended
	}

	whiteSum := 0
	blackSum := 0
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if board[i] == BLUE {
			blackSum++
		} else if board[i] == RED {
			whiteSum++
		}
	}

	player_moves := numLegalActions(board, nextToMove)
	other_moves := numLegalActions(board, opponent(nextToMove))

	if player_moves == 0 && other_moves == 0 {
		if blackSum > whiteSum {
			return gomcts.GameResult(BLUE), true
		} else if blackSum < whiteSum {
			return gomcts.GameResult(RED), true
		}
		// draw
		return gomcts.GameResult(EMPTY), true
	}

	return gomcts.GameResult(EMPTY), false
}

// ApplyTo - OthelloBoardGameAction implementation of ApplyTo method of Action interface
func (a OthelloBoardGameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	g := s.(OthelloGameState)
	board := make([]int, BOARD_SIZE)

	copy(board, g.board)
	g.board = board

	if g.nextToMove != a.value {
		panic("*hands slapped*,  not your turn")
	}

	if g.board[a.move] != EMPTY {
		panic("*hands slapped*,  square already occupied")
	}

	if !bound(a.move) {
		panic("*hands slapped*,  move out of bounds")
	}

	makeMove(a.move, g)

	g.nextToMove = opponent(a.value)

	// Next to play has no moves
	if numLegalActions(g.board, g.nextToMove) == 0 {
		g.nextToMove = a.value
	}

	return g
}

// GetLegalActions - OthelloGameState implementation of GetLegalAction method of GameState interface
func (s OthelloGameState) GetLegalActions() []gomcts.Action {
	cnt := 0
	actions := make([]gomcts.Action, 0, 0)
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if legalMove(board, i, nextToMove) {
			cnt++
			actions = append(actions, OthelloBoardGameAction{move: i, value: s.nextToMove})
		}
	}

	return actions
}

// NextToMove - OthelloGameState implementation of NextToMove method of GameState interface
func (s OthelloGameState) NextToMove() int {
	return s.nextToMove
}

/*
 * OthelloBoardGameState custom methods
 */

// GetBoard - Get a copy of current board
func (s OthelloGameState) GetBoard() []int {
	board := make([]int, len(s.board))
	copy(board, s.board)
	return board
}

// Clone - Clone the current game stat
func (s OthelloGameState) Clone() OthelloGameState {
	board := make([]int, BOARD_SIZE)
	copy(board, s.board)
	state := OthelloGameState{nextToMove: s.nextToMove, board: board}
	return state
}

// GetScore - Get the current score
func (s OthelloGameState) GetScore() (p1, p2 int) {
	return count(s.board, BLUE), count(s.board, RED)
}

/*
 * OthelloBoardGameAction custom methods
 */

// GetMove - Get move field of
func (a OthelloBoardGameAction) GetMove() int {
	return a.move
}

// GetValue
func (a OthelloBoardGameAction) GetValue() int {
	return a.value
}
