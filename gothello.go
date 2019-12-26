package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/unathi-skosana/gothello/pkg/gomcts"
	"github.com/unathi-skosana/gothello/pkg/othello"
)

const depth = 1000
const BOARD_SIZE = othello.BOARD_WIDTH
const board_row_top = "+---+---+---+---+---+---+---+---+"

var st = tcell.StyleDefault
var clock *Clock

// DIR
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
	i, j := 4, 5
	s, e := tcell.NewScreen()

	eval, player, level := parsArgs()

	// BLUE always goes first.
	var gs gomcts.GameState = othello.New(othello.BLUE)
	clock = newClock()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	quit := make(chan struct{})

	var refresh = func() {
		s.Clear()
		printGame(s, gs.(othello.OthelloGameState), level, gs.NextToMove() == player, i, j)
	}

	clock.TickFunc = refresh

	go func() {
		for {
			// blocking
			if gs.NextToMove() != player {
				action := gomcts.MonteCarloTreeSearch(gs, eval, depth)
				gs = action.ApplyTo(gs)
				refresh()
			}
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRight: // right
					i, j = moveSelector(E, i, j)
				case tcell.KeyLeft: // left
					i, j = moveSelector(W, i, j)
				case tcell.KeyDown: // down
					i, j = moveSelector(S, i, j)
				case tcell.KeyUp: // up
					i, j = moveSelector(N, i, j)
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
						i, j = moveSelector(W, i, j)
					case 108: // l
						i, j = moveSelector(E, i, j)
					case 106: // j
						i, j = moveSelector(S, i, j)
					case 107: // k
						i, j = moveSelector(N, i, j)
					case 110: // n
						gs = othello.New(othello.BLUE)
						clock = newClock()
						clock.TickFunc = refresh
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

type Clock struct {
	ticker   *time.Ticker
	Tick     bool
	Duration time.Duration
	TickFunc func()
}

func newClock() *Clock {
	clock := &Clock{
		ticker: time.NewTicker(time.Millisecond * 500),
		Tick:   true,
	}
	t0 := time.Now()

	go func() {
		for t := range clock.ticker.C {
			clock.Tick = !clock.Tick
			clock.Duration = t.Sub(t0)
			if clock.TickFunc != nil {
				clock.TickFunc()
			}
		}
	}()

	return clock
}

// parse and process arguments
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
			eval = othello.OthelloRandomRolloutPolicy
			level = "easy"
		default:
			panic("Invalid argument for -d flag. See help")
		}

		switch *player {
		case "red":
			nextToMove = 2
		case "blue":
			nextToMove = 1
		default:
			panic("Invalid argument for -p flag. See help")
		}
	}

	return eval, nextToMove, level

}

// check if new bounded position
func moveSelector(d, i, j int) (int, int) {
	_i, _j := nxt(d, i, j)
	if bound(_i) && bound(_j) {
		return _i, _j
	} else {
		return i, j
	}
}

// given a direction to move to and current position get new position
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

// checks if point is bounded by the board size
func bound(i int) bool {
	return i >= 0 && i < BOARD_SIZE
}

// prints string on screen
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

// prints current game state
func printGame(s tcell.Screen, gs othello.OthelloGameState, level string, showLegalMoves bool, ci, cj int) {
	const header = 3
	const w = tcell.ColorWhite
	const b = tcell.ColorBlue
	const r = tcell.ColorRed

	nextToMove := gs.NextToMove()
	board := gs.GetBoard()
	actions := gs.GetLegalActions()

	XOFF := 10
	YOFF := 5
	SYMBOLS := []string{" ", "•", "•", "x"}
	COLORS := []tcell.Color{w, b, r, w}

	// board numbers
	for i := 0; i < BOARD_SIZE; i++ {
		n := fmt.Sprintf("%d", i)
		puts(s, w, XOFF+i*4+2, YOFF+header-1, n)
		puts(s, w, XOFF+i*4+2, YOFF+2*BOARD_SIZE+header+1, n)
		puts(s, w, XOFF-header, YOFF+i*2+header+1, n)
		puts(s, w, XOFF+BOARD_SIZE*4+2, YOFF+i*2+header+1, n)
	}

	// game state
	for i := 0; i < BOARD_SIZE; i++ {
		puts(s, w, XOFF, YOFF+header+2*i, board_row_top)
		for j := 0; j < BOARD_SIZE+1; j++ {
			piece := board[i+1+10*(j+1)]
			if piece == othello.BLUE || piece == othello.RED {
				puts(s, COLORS[piece], XOFF+2*i*2+2, YOFF+2*j+header+1, SYMBOLS[piece])
			}
			puts(s, w, XOFF+4*j, YOFF+header+2*i+1, "|")
		}
	}
	puts(s, w, XOFF, YOFF+header+2*BOARD_SIZE, board_row_top)

	// user control player is not playing do not show selector and possible
	// actions
	if showLegalMoves {
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
	}

	// level
	puts(s, w, XOFF+BOARD_SIZE*4+16+2, YOFF+header-1, fmt.Sprintf("level: %v", level))

	// controls
	puts(s, w, XOFF+BOARD_SIZE*4+16+2, YOFF+header+1, "Movement")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+2, "h - Move left")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+3, "j - Move down")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+4, "k - Move up")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+5, "l - Move right")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+6, "Enter/Space - Place piece")

	// commands
	puts(s, w, XOFF+BOARD_SIZE*4+16+2, YOFF+header+8, "Commands")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+9, "q - Quit")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+10, "n - New game")
	puts(s, w, XOFF+BOARD_SIZE*4+16+3, YOFF+header+11, "r - Reload")

	// score
	p1, p2 := gs.GetScore()
	puts(s, w, XOFF+BOARD_SIZE*2-5, YOFF, fmt.Sprintf("_ %2d - %-2d _", p1, p2))
	puts(s, COLORS[othello.BLUE], XOFF+BOARD_SIZE*2-5, YOFF, SYMBOLS[othello.BLUE])
	puts(s, COLORS[othello.RED], XOFF+BOARD_SIZE*2+5, YOFF, SYMBOLS[othello.RED])

	// time
	mins := int(clock.Duration.Minutes())
	secs := int(clock.Duration.Seconds()) % 60
	deli := ":"
	if !clock.Tick {
		deli = " "
	}
	time := fmt.Sprintf("%02d%s%02d", mins, deli, secs)
	puts(s, w, XOFF+BOARD_SIZE*2-len(time)/2, YOFF+header+2*BOARD_SIZE+3, time)

	s.Show()
}
