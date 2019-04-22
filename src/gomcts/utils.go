package gomcts

const EMPTY = 0
const BLACK = 1
const WHITE = -1
const OUTER = 2
const BOARDSIZE = 100
const REALBOARDSIZE = 64

var ALLDIRECTIONS = []int8{-11, -10, -9, -1, 1, 9, 10, 11}

//  Heuritic Weights
const PARITY_WEIGHT = 250.00
const MOBILITY_WEIGHT = 39.22
const CORNERS_WEIGHT = 801.724
const FRONTIERS_WEIGHT = 74.396

func initializeOthelloBoard() []int8 {
	board := make([]int8, BOARDSIZE)
	for i := 0; i <= 9; i++ {
		board[i] = OUTER
	}
	for i := 10; i <= 89; i++ {
		if i%10 >= 1 && i%10 <= 8 {
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

func copyOthelloBoard(board []int8) []int8 {
	newboard := make([]int8, BOARDSIZE)
	for i := 0; i < BOARDSIZE; i++ {
		newboard[i] = board[i]
	}
	return newboard
}

func copyOthelloGameState(s OthelloGameState) OthelloGameState {
	board := copyOthelloBoard(s.board)
	state := OthelloGameState{nextToMove: s.nextToMove, board: board, emptySquares: s.emptySquares}
	return state
}

func evalFunc(s OthelloGameState) float64 {

	var parityHeuristic float64
	var mobilityHeuristic float64
	var cornersHeuristic float64
	var frontiersHeuristic float64
	var playerFrontiers float64
	var oppFrontiers float64
	var playerPieces float64
	var oppPieces float64
	var playerCorners float64
	var oppCorners float64
	var playerMoves float64
	var oppMoves float64

	var col int8
	var row int8

	player := s.nextToMove
	board := s.board
	playerPieces = float64(count(player, board))
	oppPieces = float64(count(opponent(player), board))

	// Frontiers

	for row = 1; row <= 8; row++ {
		for col = 1; col <= 8; col++ {
			if board[col+(10*row)] != EMPTY {
				for dir := 0; dir < 8; dir++ {
					dirMove := col + (10 * row) + ALLDIRECTIONS[dir]
					if dirMove >= 11 && dirMove <= 88 && board[dirMove] == EMPTY {
						if board[dirMove+ALLDIRECTIONS[dir]] == player {
							playerFrontiers++
						} else {
							oppFrontiers++
						}
					}
				}
			}
		}
	}

	if playerFrontiers > oppFrontiers {
		frontiersHeuristic =
			-100 * playerFrontiers / (playerFrontiers + oppFrontiers)
	} else if oppFrontiers > playerFrontiers {
		frontiersHeuristic =
			100 * oppFrontiers / (playerFrontiers + oppFrontiers)
	} else {
		frontiersHeuristic = 0.00
	}

	if playerPieces > oppPieces {
		parityHeuristic = 100 * playerPieces / (playerPieces + oppPieces)
	} else if oppPieces > playerPieces {
		parityHeuristic = -100 * oppPieces / (playerPieces + oppPieces)
	} else {
		parityHeuristic = 0.00
	}

	// Mobility

	playerMoves = float64(numLegalActions(s))
	s.nextToMove *= -1
	oppMoves = float64(numLegalActions(s))
	s.nextToMove *= -1

	if playerMoves > oppMoves {
		mobilityHeuristic = 100 * (playerMoves) / (playerMoves + oppMoves)
	} else if oppMoves > playerMoves {
		mobilityHeuristic = -100 * (oppMoves) / (playerMoves + oppMoves)
	} else {
		mobilityHeuristic = 0.00
	}

	if board[1+(10*8)] == player {
		playerCorners++
	} else if board[1+(10*8)] == opponent(player) {
		oppCorners++
	}

	if board[1+(10*1)] == player {
		playerCorners++
	} else if board[1+(10*1)] == opponent(player) {
		oppCorners++
	}

	if board[8+(10*1)] == player {
		playerCorners++
	} else if board[8+(10*1)] == opponent(player) {
		oppCorners++
	}

	if board[8+(10*8)] == player {
		playerCorners++
	} else if board[8+(10*8)] == opponent(player) {
		oppCorners++
	}

	cornersHeuristic = playerCorners - oppCorners

	// final score
	return (PARITY_WEIGHT * parityHeuristic) +
		(MOBILITY_WEIGHT * mobilityHeuristic) +
		(CORNERS_WEIGHT * cornersHeuristic) +
		(FRONTIERS_WEIGHT * frontiersHeuristic)
}

func makeMove(move int8, s OthelloGameState) {
	s.board[move] = s.nextToMove
	for i := 0; i <= 7; i++ {
		makeFlips(move, ALLDIRECTIONS[i], s)
	}
}

func makeFlips(move int8, dir int8, s OthelloGameState) {
	var bracketer int8
	var c int8

	bracketer = wouldFlip(move, dir, s)

	if bracketer != 0 {
		c = move + dir
		for {
			s.board[c] = s.nextToMove
			c = c + dir
			if c == bracketer {
				break
			}
		}
	}
}

func getArrIndex(s int8) int8 {
	var row int8
	var col int8

	row = s / 10
	col = s % 10

	return (10 * (row + 1)) + col + 1
}

func numLegalActions(s OthelloGameState) int8 {
	var cnt int8
	for i := 11; i <= 88; i++ {
		if legalPlayer(int8(i), s) == 1 {
			cnt++
		}
	}
	return cnt
}

func legalPlayer(move int8, s OthelloGameState) int8 {
	var i int8

	if validPlayer(move) == 0 {
		return 0
	}

	if s.board[move] == EMPTY {
		i = 0
		for {
			if !(i <= 7 && wouldFlip(move, ALLDIRECTIONS[i], s) == 0) {
				break
			}
			i++
		}

		if i == 8 {
			return 0
		}
		return 1

	}

	return 0
}

func validPlayer(move int8) int8 {
	if (move >= 11) && (move <= 88) && (move%10 >= 1) && (move%10 <= 8) {
		return 1
	}
	return 0
}

func wouldFlip(move int8, dir int8, s OthelloGameState) int8 {
	var c int8

	c = move + dir

	if s.board[c] == opponent(s.nextToMove) {
		return findBracketingPiece(c+dir, dir, s)
	}

	return 0

}

func findBracketingPiece(square int8, dir int8, s OthelloGameState) int8 {
	for s.board[square] == opponent(s.nextToMove) {
		square = square + dir
	}

	if s.board[square] == s.nextToMove {
		return square
	}

	return 0
}

func opponent(nextToMove int8) int8 {
	return nextToMove * -1
}

func nameof(piece int8) string {
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

func count(player int8, board []int8) int8 {
	var cnt int8

	for i := 1; i <= 88; i++ {
		if board[i] == player {
			cnt++
		}
	}
	return cnt
}
