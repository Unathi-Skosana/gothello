package gomcts

import "fmt"

const EMPTY = 0
const BLACK = 1
const WHITE = -1
const OUTER = 2
const BOARDSIZE = 100
const REALBOARDSIZE =  64
var ALLDIRECTIONS = []int8{-11, -10, -9, -1, 1, 9, 10, 11}

//  Heuritic Weights
const PARITY_WEIGHT = 250.00;
const MOBILITY_WEIGHT = 39.22;
const CORNERS_WEIGHT = 801.724;
const FRONTIERS_WEIGHT = 74.396;



func initializeOthelloBoard() []int8 {
    board := make([]int8, BOARDSIZE)
    for i := 0; i <= 9; i++ { board[i] = OUTER }
    for i := 10; i <= 89; i++ {
      if (i % 10 >= 1 && i % 10 <= 8) {
        board[i] = EMPTY
      } else {
        board[i] = OUTER
      }
    }

    for i := 90; i <= 99; i++ { board[i] = OUTER }

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



func CreateOthelloInitialGameState() OthelloGameState {
    board := initializeOthelloBoard()
    state := OthelloGameState{nextToMove: BLACK, board: board, emptySquares: uint16(REALBOARDSIZE)- 4}
    return state
}



func evalFunc(s OthelloGameState) float64 {

  var parity_heuristic float64 = 0.0
  var mobility_heuristic float64 = 0.0
  var corners_heuristic float64 = 0.0
  var frontiers_heuristic float64 = 0.0

  player := s.nextToMove
  board := s.board

  // Frontiers
  player_frontiers := 0.0
  opp_frontiers := 0.0

  var col int8
  var row int8

  for row = 1; row <= 8; row++ {
    for col = 1; col <= 8; col++ {
      if board[col + (10 * row)] != EMPTY {
        for dir := 0; dir < 8; dir++ {
          dir_move := col + (10 * row) + ALLDIRECTIONS[dir]
          if dir_move >= 11 && dir_move <= 88 && board[dir_move] == EMPTY {
            if board[dir_move + ALLDIRECTIONS[dir]] == player {
              player_frontiers++
            } else {
              opp_frontiers++
            }
          }
        }
      }
    }
  }

  if player_frontiers > opp_frontiers {
    frontiers_heuristic =
        -100 * player_frontiers / (player_frontiers + opp_frontiers);
  } else if (opp_frontiers > player_frontiers) {
    frontiers_heuristic =
        100 * opp_frontiers / (player_frontiers + opp_frontiers)
  } else {
    frontiers_heuristic = 0.00
  }

  // Parity
  player_pieces := float64(count(player, board))
  opp_pieces := float64(count(opponent(player), board))

  if (player_pieces > opp_pieces) {
    parity_heuristic = 100 * player_pieces / (player_pieces + opp_pieces)
  } else if (opp_pieces > player_pieces) {
    parity_heuristic = -100 * opp_pieces / (player_pieces + opp_pieces)
  } else {
    parity_heuristic = 0.00
  }

  // Mobility

  player_moves  := float64(numLegalActions(s))
  s.nextToMove *= -1
  opp_moves := float64(numLegalActions(s))
  s.nextToMove *= -1

  if (player_moves > opp_moves) {
    mobility_heuristic = 100 * (player_moves) / (player_moves + opp_moves)
  } else if (opp_moves > player_moves) {
    mobility_heuristic = -100 * (opp_moves) / (player_moves + opp_moves)
  } else {
    mobility_heuristic = 0.00
  }

  // corners_heuristic
  player_corners := 0.0
  opp_corners := 0.0

  if (board[1 + (10 * 8)] == player) {
    player_corners++
  } else if (board[1 + (10 * 8)] == opponent(player)) {
    opp_corners++
  }

  if (board[1 + (10 * 1)] == player) {
    player_corners++
  } else if (board[1 + (10 * 1)] == opponent(player)) {
    opp_corners++
  }

  if (board[8 + (10 * 1)] == player) {
    player_corners++
  } else if (board[8 + (10 * 1)] == opponent(player)) {
    opp_corners++
  }

  if (board[8 + (10 * 8)] == player) {
    player_corners++
  } else if (board[8 + (10 * 8)] == opponent(player)) {
    opp_corners++
  }

  corners_heuristic = player_corners - opp_corners

  // final score
  return (PARITY_WEIGHT * parity_heuristic) +
         (MOBILITY_WEIGHT * mobility_heuristic) +
         (CORNERS_WEIGHT * corners_heuristic) +
         (FRONTIERS_WEIGHT * frontiers_heuristic);
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

  if (bracketer != 0) {
    c = move + dir
    for {
      s.board[c] = s.nextToMove
      c = c + dir
      if (c == bracketer) {
        break;
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
  var cnt int8 = 0
  for i := 11; i <= 88; i++ {
    if legalPlayer(int8(i), s) == 1 {
      cnt ++
    }
  }
  return cnt
}


func legalPlayer(move int8, s OthelloGameState) int8 {
  var i int8

  if (validPlayer(move) == 0) {
    return 0
  }

  if (s.board[move] == EMPTY) {
    i = 0
    for {
      if !((i <= 7 && wouldFlip(move, ALLDIRECTIONS[i], s) == 0)) { break; }
        i++;
    }

    if i == 8 {
      return 0
    }
    return 1

  }

  return 0
}

func validPlayer(move int8) int8 {
  if ((move >= 11) && (move <= 88) && (move % 10 >= 1) && (move % 10 <= 8)) {
    return 1
  }
  return 0
}

func wouldFlip(move int8, dir int8, s OthelloGameState) int8 {
  var c int8

  c = move + dir

  if (s.board[c] == opponent(s.nextToMove)) {
    return findBracketingPiece(c + dir, dir, s)
  }

  return 0

}

func findBracketingPiece(square int8, dir int8, s OthelloGameState) int8 {
  for s.board[square] == opponent(s.nextToMove) {
    square = square + dir
  }

  if (s.board[square] == s.nextToMove) {
    return square
  }

  return 0
}

func opponent(nextToMove int8) int8 {
  return nextToMove * -1
}

func PrintBoard(s OthelloGameState) {
  fmt.Printf("   1 2 3 4 5 6 7 8 [%s=%d %s=%d]\n", nameof(BLACK), count(BLACK, s.board), nameof(WHITE), count(WHITE, s.board))
  for row := 1; row <= 8; row++ {
    fmt.Printf("%d  ", row)
    for col := 1; col <= 8; col++ {
      fmt.Printf("%s ", nameof(s.board[col + (10 * row)]))
    }
    fmt.Printf("\n")
  }
}

func nameof(piece int8) string {
  if (piece == EMPTY ) { return "." }
  if (piece == BLACK) { return "b" }
  if (piece == WHITE) { return "w" }

  return "?"

}

func count(player int8, board []int8) int8 {
  var cnt int8 = 0

  for i := 1; i <= 88; i++ {
    if (board[i] == player) {
      cnt++
    }
  }
  return cnt
}
