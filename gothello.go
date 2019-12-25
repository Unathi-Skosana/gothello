package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/unathi-skosana/gothello/pkg/gomcts"
	"github.com/unathi-skosana/gothello/pkg/othello"
)

const BOARD_SIZE = othello.BOARD_WIDTH
const board_row_top = "+---+---+---+---+---+---+---+---+"

var st = tcell.StyleDefault

/*
package main

import (
	"fmt"

	"github.com/unathi-skosana/gothello/othello"
)

func main() {

	var i int = 1
	var chosenAction gomcts.Action
	var s gomcts.GameState = othello.New()

	othello.PrintBoard(s)

	for !s.IsGameEnded() {
		if i%2 == 1 {
			chosenAction = gomcts.MonteCarloTreeSearch(s, othello.OthelloMediumRolloutPolicy, 1500)
		} else {
			chosenAction = gomcts.MonteCarloTreeSearch(s, othello.OthelloHardRolloutPolicy, 1500)
		}
		s = chosenAction.ApplyTo(s)
		othello.PrintBoard(s)
		fmt.Println(chosenAction)
		i++
	}
}

*/

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

func main() {
	i, j := 0, 0
	s, e := tcell.NewScreen()

	var gs gomcts.GameState = othello.New()

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
		printGame(s, gs.(othello.OthelloGameState), i, j)
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
				case tcell.KeyRune:
					key := ev.Rune()
					switch key {
					case 32: // Space
					case 104: // h
						i, j = nxtBound(W, i, j)
					case 108: // l
						i, j = nxtBound(E, i, j)
					case 106: // j
						i, j = nxtBound(S, i, j)
					case 107: // k
						i, j = nxtBound(N, i, j)
					case 110: // n
						// initialise new game
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

func printGame(s tcell.Screen, gs othello.OthelloGameState, ci, cj int) {
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
	SYMBOLS := []string{" ", "●", "●", "+"}
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
			puts(s, COLORS[piece], XOFF+2*i*2+2, YOFF+2*j+header+1, SYMBOLS[piece])
			puts(s, d, XOFF+4*j, YOFF+header+2*i+1, "|")
		}
	}
	puts(s, d, XOFF, YOFF+header+2*BOARD_SIZE, board_row_top)

	// Legal actions for current player
	for i := 0; i < len(actions); i++ {
		action := actions[i].(othello.OthelloBoardGameAction)
		move := action.GetMove()
		value := action.GetValue()
		x := move/10 - 1
		y := move%10 - 1

		puts(s, COLORS[value], XOFF+2*x*2+2, YOFF+2*y+header+1, SYMBOLS[3])

	}

	// selector
	puts(s, COLORS[nextToMove], XOFF+2*ci*2+1, YOFF+2*cj+header+1, "[")
	puts(s, COLORS[nextToMove], XOFF+2*ci*2+3, YOFF+2*cj+header+1, "]")

	// control
	puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header+1, "Movement")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+2, "h - Move left")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+3, "j - Move down")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+4, "k - Move up")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+5, "l - Move right")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+6, "Enter/Space - Place piece")

	puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header+8, "Commands")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+9, "q - Quit")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+10, "n - New game")
	puts(s, d, XOFF+BOARD_SIZE*4+16+3, YOFF+header+11, "r - Reload")

	// score
	puts(s, d, XOFF+BOARD_SIZE*2-5, YOFF, fmt.Sprintf("_ %2d - %-2d _", 0, 0))
	puts(s, COLORS[1], XOFF+BOARD_SIZE*2-5, YOFF, SYMBOLS[1])
	puts(s, COLORS[2], XOFF+BOARD_SIZE*2+5, YOFF, SYMBOLS[2])

}
