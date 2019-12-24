package gomcts

import (
	"errors"
	"fmt"
)

// Pieces
const (
	EMPTY = iota
	BLACK
	WHITE
	OUTER
)

// DIR
const DIR = 8
const (
	E = iota
	NE
	N
	NW
	W
	SW
	S
	SE
)

// sizes
const BOARD_SIZE = 100
const BOARD_WIDTH = 8
const PIECE_SLOTS = 64

// BLOCKS
const FIRST_BLOCK = 11
const LAST_BLOCK = 88

// NO MORE MOVES
const NO_MOVES = -1

// DIRECTIONS
var ALLDIRECTIONS = []int{-11, -10, -9, -1, 1, 9, 10, 11}

// OthelloBoardGameAction - action on a othello board
type OthelloBoardGameAction struct {
	move  int
	value int
}

// OthelloGameState - Othello game state
type OthelloGameState struct {
	nextToMove   int
	board        []int
	emptySquares uint16
	ended        bool
	result       GameResult
}

// PrintBoard - prints out game board to console
func (s OthelloGameState) PrintBoard() {
	board := s.board
	fmt.Printf("   1 2 3 4 5 6 7 8 [%s=%d %s=%d]\n", nameof(BLACK), count(board, BLACK), nameof(WHITE), count(board, WHITE))
	for row := 1; row <= BOARD_WIDTH; row++ {
		fmt.Printf("%d  ", row)
		for col := 1; col <= BOARD_WIDTH; col++ {
			fmt.Printf("%s ", nameof(board[col+(10*row)]))
		}
		fmt.Printf("\n")
	}
}

// CreateOthelloInitialGameState - initializes a othello game state
func CreateOthelloInitialGameState() OthelloGameState {
	board := initializeOthelloBoard()
	state := OthelloGameState{nextToMove: BLACK, board: board, emptySquares: uint16(PIECE_SLOTS) - 4}
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
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if board[i] == BLACK {
			blackSum++
		} else if board[i] == WHITE {
			whiteSum++
		}
	}

	player_moves := numLegalActions(board, nextToMove)
	other_moves := numLegalActions(board, opponent(nextToMove))

	if s.emptySquares == 0 || player_moves == 0 && other_moves == 0 {
		if blackSum > whiteSum {
			return GameResult(BLACK), true
		} else if blackSum < whiteSum {
			return GameResult(WHITE), true
		}
		return GameResult(0), true
	}

	return GameResult(0), false
}

func (a OthelloBoardGameAction) GetMove() int {
	return a.move
}

// ApplyTo - OthelloBoardGameAction implementation of ApplyTo method of Action interface
func (a OthelloBoardGameAction) ApplyTo(s GameState) GameState {
	OthelloGameState := s.(OthelloGameState)
	OthelloGameState.board = copyOthelloBoard(OthelloGameState.board)

	if a.move == NO_MOVES {
		OthelloGameState.nextToMove = opponent(OthelloGameState.nextToMove)
		return OthelloGameState
	}

	if OthelloGameState.nextToMove != a.value {
		panic("*hands slapped*,  not your turn")
	}

	makeMove(a.move, OthelloGameState)

	OthelloGameState.nextToMove = opponent(OthelloGameState.nextToMove)

	OthelloGameState.emptySquares--

	return OthelloGameState
}

// GetLegalActions - OthelloGameState implementation of GetLegalAction method of GameState interface
func (s OthelloGameState) GetLegalActions() []Action {
	cnt := 0
	actions := make([]Action, 0, 0)
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if legalPlayer(board, i, nextToMove) == 1 {
			cnt++
			actions = append(actions, OthelloBoardGameAction{move: i, value: s.nextToMove})
		}
	}

	if cnt == 0 {
		actions = append(actions, OthelloBoardGameAction{move: NO_MOVES, value: s.nextToMove})
	}

	return actions
}

// NextToMove - OthelloGameState implementation of NextToMove method of GameState interface
func (s OthelloGameState) NextToMove() int {
	return s.nextToMove
}

// NextBound - ...
func NextBound(d, i, j int) (int, int) {
	_i, _j := nxt(d, i, j)
	if bound(_i) && bound(_j) {
		return _i, _j
	} else {
		return i, j
	}
}

// Utils
func nxt(d, i, j int) (int, int) {
	switch d {
	case E:
		return i + 1, j
	case NE:
		return i + 1, j - 1
	case N:
		return i, j - 1
	case NW:
		return i - 1, j - 1
	case W:
		return i - 1, j
	case SW:
		return i - 1, j + 1
	case S:
		return i, j + 1
	case SE:
		return i + 1, j + 1
	default:
		panic(errors.New("unknown direction"))
	}
}

func bound(i int) bool {
	return i >= 0 && i < BOARD_WIDTH
}

func initializeOthelloBoard() []int {
	board := make([]int, BOARD_SIZE)
	for i := 0; i <= 9; i++ {
		board[i] = OUTER
	}
	for i := FIRST_BLOCK - 1; i <= LAST_BLOCK+1; i++ {
		if i%10 >= 1 && i%10 <= BOARD_WIDTH {
			board[i] = EMPTY
		} else {
			board[i] = OUTER
		}
	}

	for i := 90; i <= 99; i++ {
		board[i] = OUTER
	}

	board[44] = WHITE
	board[45] = BLACK
	board[54] = BLACK
	board[55] = WHITE

	return board

}

func copyOthelloBoard(board []int) []int {
	newboard := make([]int, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		newboard[i] = board[i]
	}
	return newboard
}

func copyOthelloGameState(s OthelloGameState) OthelloGameState {
	board := copyOthelloBoard(s.board)
	state := OthelloGameState{nextToMove: s.nextToMove, board: board, emptySquares: s.emptySquares}
	return state
}

func evalFunc(s OthelloGameState, parityWeight, mobilityWeight, cornersWeight, frontiersWeight float64) float64 {
	nextToMove := s.nextToMove
	board := s.board

	// Frontiers
	playerFrontiers := 0.0
	oppFrontiers := 0.0

	for row := 1; row <= DIR; row++ {
		for col := 1; col <= DIR; col++ {
			if board[col+(10*row)] != EMPTY {
				for dir := 0; dir < DIR; dir++ {
					dirMove := col + (10 * row) + ALLDIRECTIONS[dir]
					if dirMove >= FIRST_BLOCK && dirMove <= LAST_BLOCK && board[dirMove] == EMPTY {
						if board[dirMove+ALLDIRECTIONS[dir]] == nextToMove {
							playerFrontiers++
						} else {
							oppFrontiers++
						}
					}
				}
			}
		}
	}

	frontiersHeuristic := 0.00

	if playerFrontiers > oppFrontiers {
		frontiersHeuristic = -100.00 * playerFrontiers / (playerFrontiers + oppFrontiers)
	} else if oppFrontiers > playerFrontiers {
		frontiersHeuristic =
			100 * oppFrontiers / (playerFrontiers + oppFrontiers)
	}

	// Parity

	playerPieces := float64(count(board, nextToMove))
	oppPieces := float64(count(board, opponent(nextToMove)))

	parityHeuristic := 0.00

	if playerPieces > oppPieces {
		parityHeuristic = 100 * playerPieces / (playerPieces + oppPieces)
	} else if oppPieces > playerPieces {
		parityHeuristic = -100 * oppPieces / (playerPieces + oppPieces)
	}

	// Mobility

	playerMoves := float64(numLegalActions(board, nextToMove))
	oppMoves := float64(numLegalActions(board, opponent(nextToMove)))

	mobilityHeuristic := 0.00

	if playerMoves > oppMoves {
		mobilityHeuristic = 100 * playerMoves / (playerMoves + oppMoves)
	} else if oppMoves > playerMoves {
		mobilityHeuristic = -100 * oppMoves / (playerMoves + oppMoves)
	} else {
	}

	// Corners

	playerCorners := 0.0
	oppCorners := 0.0

	if board[1+(10*BOARD_WIDTH)] == nextToMove {
		playerCorners++
	} else if board[1+(10*BOARD_WIDTH)] == opponent(nextToMove) {
		oppCorners++
	}

	if board[1+(10*1)] == nextToMove {
		playerCorners++
	} else if board[1+(10*1)] == opponent(nextToMove) {
		oppCorners++
	}

	if board[BOARD_WIDTH+(10*1)] == nextToMove {
		playerCorners++
	} else if board[BOARD_WIDTH+(10*1)] == opponent(nextToMove) {
		oppCorners++
	}

	if board[BOARD_WIDTH+(10*BOARD_WIDTH)] == nextToMove {
		playerCorners++
	} else if board[BOARD_WIDTH+(10*BOARD_WIDTH)] == opponent(nextToMove) {
		oppCorners++
	}

	cornersHeuristic := playerCorners - oppCorners

	// final score
	return parityWeight*parityHeuristic +
		mobilityWeight*mobilityHeuristic +
		cornersWeight*cornersHeuristic +
		frontiersWeight*frontiersHeuristic
}

func makeMove(move int, s OthelloGameState) {
	nextToMove := s.nextToMove
	board := s.board

	board[move] = nextToMove

	for i := 0; i < DIR; i++ {
		makeFlips(board, move, ALLDIRECTIONS[i], nextToMove)
	}
}

func makeFlips(board []int, move, dir, nextToMove int) {
	var bracketer int
	var c int

	bracketer = wouldFlip(board, move, dir, nextToMove)

	if bracketer != 0 {
		c = move + dir
		for {
			board[c] = nextToMove
			c = c + dir
			if c == bracketer {
				break
			}
		}
	}
}

func idx(i, j int) int {
	return (10 * (i + 1)) + j + 1
}

func numLegalActions(board []int, nextToMove int) int {
	var cnt int
	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if legalPlayer(board, i, nextToMove) == 1 {
			cnt++
		}
	}
	return cnt
}

func legalPlayer(board []int, move, nextToMove int) int {
	var i int

	if validPlayer(move) == 0 {
		return 0
	}

	if board[move] == EMPTY {
		i = 0
		for {
			if !(i < DIR && wouldFlip(board, move, ALLDIRECTIONS[i], nextToMove) == 0) {
				break
			}
			i++
		}

		if i == DIR {
			return 0
		}

		return 1

	}

	return 0
}

func validPlayer(move int) int {
	if (move >= FIRST_BLOCK) && (move <= LAST_BLOCK) && (move%10 >= 1) && (move%10 <= BOARD_WIDTH) {
		return 1
	}
	return 0
}

func wouldFlip(board []int, move, dir, nextToMove int) int {
	var c int

	c = move + dir

	if board[c] == opponent(nextToMove) {
		return findBracketingPiece(board, c+dir, dir, nextToMove)
	}

	return 0

}

func findBracketingPiece(board []int, square, dir, nextToMove int) int {
	for board[square] == opponent(nextToMove) {
		square = square + dir
	}

	if board[square] == nextToMove {
		return square
	}

	return 0
}

func opponent(nextToMove int) int {
	if nextToMove == BLACK {
		return WHITE
	} else if nextToMove == WHITE {
		return BLACK
	}

	panic("*hands slapped*,  invalid player")

}

func nameof(piece int) string {
	if piece == EMPTY {
		return "."
	}
	if piece == BLACK {
		return "b"
	}
	if piece == WHITE {
		return "w"
	}

	return "?"

}

func count(board []int, nextToMove int) int {
	var cnt int

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if board[i] == nextToMove {
			cnt++
		}
	}
	return cnt
}
