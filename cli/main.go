package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

const BOARD_SIZE = 8
const board_row_top = "+---+---+---+---+---+---+---+---+"
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

var st = tcell.StyleDefault

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

func nextBound(d, i, j int) (int, int) {
	_i, _j := nxt(d, i, j)
	if bound(_i) && bound(_j) {
		return _i, _j
	} else {
		return i, j
	}
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

func printGame(s tcell.Screen, ci, cj int) {
	const header = 3
	const d = tcell.ColorBlack
	const b = tcell.ColorBlue
	const m = tcell.ColorPurple
	const g = tcell.ColorGreen
	const y = tcell.ColorYellow
	const r = tcell.ColorRed

	XOFF := 10
	YOFF := 5
	SYMBOLS := []string{" ", "●", "●", "+", "+"}
	COLORS := []tcell.Color{d, b, r, d, d}

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
		if i == 2 {
			puts(s, d, XOFF+BOARD_SIZE*4+16+2, YOFF+header+2*i+1, "SCORE")
		}

		if i == 3 {
			puts(s, d, XOFF+BOARD_SIZE*4+16, YOFF+header+2*i+1, fmt.Sprintf("_ %2d - %-2d _", 0, 0))
			puts(s, COLORS[1], XOFF+BOARD_SIZE*4+16, YOFF+header+2*i+1, SYMBOLS[1])
			puts(s, COLORS[2], XOFF+BOARD_SIZE*4+26, YOFF+header+2*i+1, SYMBOLS[2])
		}

		for j := 0; j < BOARD_SIZE+1; j++ {
			puts(s, d, XOFF+4*j, YOFF+header+2*i+1, "|")
		}
	}
	puts(s, d, XOFF, YOFF+header+2*BOARD_SIZE, board_row_top)

	// selector
	puts(s, COLORS[2], XOFF+2*ci*2+1, YOFF+2*cj+header+1, "[")
	puts(s, COLORS[2], XOFF+2*ci*2+3, YOFF+2*cj+header+1, "]")

}

func main() {
	i, j := 7, 7

	s, e := tcell.NewScreen()
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
		printGame(s, i, j)
		s.Show()
	}
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyRight:
					i, j = nextBound(E, i, j)
				case tcell.KeyLeft:
					i, j = nextBound(W, i, j)
				case tcell.KeyDown:
					i, j = nextBound(S, i, j)
				case tcell.KeyUp:
					i, j = nextBound(N, i, j)
				case tcell.KeyRune:
					switch ev.Rune() {
					case 113:
						close(quit)
						return
					case 104:
						i, j = nextBound(W, i, j)
					case 108:
						i, j = nextBound(E, i, j)
					case 106:
						i, j = nextBound(S, i, j)
					case 107:
						i, j = nextBound(N, i, j)
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
