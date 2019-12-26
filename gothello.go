package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/unathi-skosana/gothello/pkg/gomcts"
	"github.com/unathi-skosana/gothello/pkg/othello"
)

const BOARD_SIZE = othello.BOARD_WIDTH
const board_row_top = "+---+---+---+---+---+---+---+---+"

var st = tcell.StyleDefault

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

func parsArgs() (d gomcts.RolloutPolicy, p int, l string) {
	// Create new parser object
	parser := argparse.NewParser("gothello", "")

	// player ~ blue always starts
	player := parser.String("p", "player", &argparse.Options{Required: false, Help: "Choose between : blue, red"})

	// difficulty
	difficulty := parser.String("d", "difficulty", &argparse.Options{Required: false, Help: "Choose between: easy, medium, hard"})

	eval := othello.OthelloRandomRolloutPolicy
	level := "easy"
	nextToMove := 1

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	} else {
		switch *difficulty {
		case "medium":
			eval = othello.OthelloMediumRolloutPolicy
			level = "medium"
		case "hard":
			eval = othello.OthelloHardRolloutPolicy
			level = "hard"
		case "easy":

		default:
			eval = othello.OthelloRandomRolloutPolicy
			level = "easy"
		}

		switch *player {
		case "red":
			nextToMove = 2
		case "blue":
		default:

			nextToMove = 1
		}
	}

	return eval, nextToMove, level

}

func main() {
	i, j := 4, 5
	s, e := tcell.NewScreen()

	eval, player, level := parsArgs()

	var gs gomcts.GameState = othello.New(2)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	st := st.Bold(true)

	s.SetStyle(st)
	s.Clear()

	quit := make(chan struct{})

	var refresh = func() {
		s.Clear()
		printGame(s, gs.(othello.OthelloGameState), level, i, j)
		s.Show()
	}
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRight: // right
					i, j = nxtBound(E, i, j)
				case tcell.KeyLeft: // left
					i, j = nxtBound(W, i, j)
				case tcell.KeyDown: // down
					i, j = nxtBound(S, i, j)
				case tcell.KeyUp: // up
					i, j = nxtBound(N, i, j)
				case tcell.KeyEnter:
					if gs.NextToMove() == player {
						actions := gs.GetLegalActions()
						for k := 0; k < len(actions); k++ {
							action := actions[k].(othello.OthelloBoardGameAction)
							move := action.GetMove()
							y := move/10 - 1
							x := move%10 - 1
							if x == i && y == j {
								gs = action.ApplyTo(gs)
							}
						}
					}
				case tcell.KeyRune:
					key := ev.Rune()
					switch key {
					case 32: // Space
						if gs.NextToMove() == player {
							actions := gs.GetLegalActions()
							for k := 0; k < len(actions); k++ {
								action := actions[k].(othello.OthelloBoardGameAction)
								move := action.GetMove()
								y := move/10 - 1
								x := move%10 - 1
								if x == i && y == j {
									gs = action.ApplyTo(gs)
								}
							}
						}
					case 104: // h
						i, j = nxtBound(W, i, j)
					case 108: // l
						i, j = nxtBound(E, i, j)
					case 106: // j
						i, j = nxtBound(S, i, j)
					case 107: // k
						i, j = nxtBound(N, i, j)
					case 110: // n
						gs = othello.New(player)
					case 113: // q
						close(quit)
						return
					case 114: //r
						s.Sync()
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
			refresh()
		}
	}()

	<-quit

	s.Fini()
}

func nxtBound(d, i, j int) (int, int) {
	_i, _j := nxt(d, i, j)
	if bound(_i) && bound(_j) {
		return _i, _j
	} else {
		return i, j
	}
}

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
	return i >= 0 && i < BOARD_SIZE
}

func puts(s tcell.Screen, fg tcell.Color, x, y int, str string) {
	st = st.Foreground(fg)
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], st)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], st)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], st)
		i += dwidth
	}
}

func printGame(s tcell.Screen, gs othello.OthelloGameState, level string, ci, cj int) {
	const header = 3
	const d = tcell.ColorBlack
	const b = tcell.ColorBlue
	const m = tcell.ColorPurple
	const g = tcell.ColorGreen
	const y = tcell.ColorYellow
	const r = tcell.ColorRed

	nextToMove := gs.NextToMove()
	board := gs.GetBoard()
	actions := gs.GetLegalActions()

	XOFF := 10
	YOFF := 5
	SYMBOLS := []string{" ", "•", "•", "+"}
	COLORS := []tcell.Color{d, b, r, d}

	//var score [2]int

	// board numbers
	for i := 0; i < BOARD_SIZE; i++ {
		n := fmt.Sprintf("%d", i)
		puts(s, d, XOFF+i*4+2, YOFF+header-1, n)
		puts(s, d, XOFF+i*4+2, YOFF+2*BOARD_SIZE+header+1, n)
		puts(s, d, XOFF-header, YOFF+i*2+header+1, n)
		puts(s, d, XOFF+BOARD_SIZE*4+2, YOFF+i*2+header+1, n)
	}

	// game state
	for i := 0; i < BOARD_SIZE; i++ {
		puts(s, d, XOFF, YOFF+header+2*i, board_row_top)
		for j := 0; j < BOARD_SIZE+1; j++ {
			piece := board[i+1+10*(j+1)]
			if piece == othello.BLUE || piece == othello.RED {
				puts(s, COLORS[piece], XOFF+2*i*2+2, YOFF+2*j+header+1, SYMBOLS[piece])
			}
			puts(s, d, XOFF+4*j, YOFF+header+2*i+1, "|")
		}
	}
	puts(s, d, XOFF, YOFF+header+2*BOARD_SIZE, board_row_top)

	// Legal actions for current player
	for i := 0; i < len(actions); i++ {
		action := actions[i].(othello.OthelloBoardGameAction)
		move := action.GetMove()
		value := action.GetValue()

		y := move/10 - 1
		x := move%10 - 1

		puts(s, COLORS[value], XOFF+2*x*2+2, YOFF+2*y+header+1, SYMBOLS[3])

	}

	// selector
	if !gs.IsGameEnded() {
		puts(s, COLORS[nextToMove], XOFF+2*ci*2+1, YOFF+2*cj+header+1, "{")
		puts(s, COLORS[nextToMove], XOFF+2*ci*2+3, YOFF+2*cj+header+1, "}")
	}

	// level and version
	puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header-1, "level:")

	// controls
	puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header+1, "Movement")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+2, "h - Move left")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+3, "j - Move down")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+4, "k - Move up")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+5, "l - Move right")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+6, "Enter/Space - Place piece")

	// commands
	puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header+8, "Commands")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+9, "q - Quit")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+10, "n - New game")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+11, "r - Reload")

	// score
	p1, p2 := gs.GetScore()
	puts(s, d, XOFF+BOARD_SIZE*2-5, YOFF, fmt.Sprintf("_ %2d - %-2d _", p1, p2))
	puts(s, COLORS[othello.BLUE], XOFF+BOARD_SIZE*2-5, YOFF, SYMBOLS[othello.BLUE])
	puts(s, COLORS[othello.RED], XOFF+BOARD_SIZE*2+5, YOFF, SYMBOLS[othello.RED])

}
