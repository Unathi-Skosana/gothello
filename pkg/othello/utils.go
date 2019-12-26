package othello

// evaluation function
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

// makes a move
func makeMove(move int, s OthelloGameState) {
	nextToMove := s.nextToMove
	board := s.board

	board[move] = nextToMove

	for i := 0; i < DIR; i++ {
		makeFlips(board, move, ALLDIRECTIONS[i], nextToMove)
	}
}

// flips pieces bracketed by opponent's piece
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

// gets the number of legal actions for player
func numLegalActions(board []int, nextToMove int) int {
	var cnt int
	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if legalMove(board, i, nextToMove) {
			cnt++
		}
	}
	return cnt
}

// checks if a play is legal to play
func legalMove(board []int, move, nextToMove int) bool {
	var i int

	if bound(move) && board[move] == EMPTY {
		i = 0
		for {
			if !(i < DIR && wouldFlip(board, move, ALLDIRECTIONS[i], nextToMove) == 0) {
				break
			}
			i++
		}

		if i == DIR {
			return false
		} else {
			return true
		}

	}

	return false
}

// checks if a move is valid
func bound(move int) bool {
	return (move >= FIRST_BLOCK) && (move <= LAST_BLOCK) && (move%10 >= 1) && (move%10 <= BOARD_WIDTH)
}

// Checks if a piece next to current player's piece would flip by checking
// if there is bracketing piece in some direction
func wouldFlip(board []int, move, dir, nextToMove int) int {
	var c int

	c = move + dir

	if board[c] == opponent(nextToMove) {
		return findBracketingPiece(board, c+dir, dir, nextToMove)
	}

	return 0

}

// finds the position of the piece that is bracketing an opponent's piece
// in some direction
func findBracketingPiece(board []int, square, dir, nextToMove int) int {
	for board[square] == opponent(nextToMove) {
		square = square + dir
	}

	if board[square] == nextToMove {
		return square
	}

	return 0
}

// gets the opponent of the player
func opponent(nextToMove int) int {
	if nextToMove == BLUE {
		return RED
	} else if nextToMove == RED {
		return BLUE
	}

	panic("*hands slapped*,  invalid player")

}

// counts the number of pieces on the board belonging to a player
func count(board []int, nextToMove int) int {
	var cnt int

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if board[i] == nextToMove {
			cnt++
		}
	}
	return cnt
}
