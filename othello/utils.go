package othello

import (
	"fmt"

	"github.com/unathi-skosana/gothello/gomcts"
)

// PrintBoard - prints out game board to console
func PrintBoard(g gomcts.GameState) {
	s := g.(OthelloGameState)
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

func evaluate(s OthelloGameState, parityWeight, mobilityWeight, cornersWeight, frontiersWeight float64) float64 {
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
		return "+"
	}
	if piece == BLACK {
		return "@"
	}
	if piece == WHITE {
		return "X"
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
