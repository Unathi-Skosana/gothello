package main

import (
	"encoding/json"
	"fmt"
	"gomcts"
	"html/template"
	"net/http"

	"github.com/yosssi/ace"
)

const START_OP = "start"
const REQUEST_OP = "reset"
const MOVE_OP = "move"
const LEGAL_MOVES = "lmoves"

const HUMAN = "1"
const RANDOM_AI = "2"
const SMART_AI = "3"

// Request from the running game
type OthelloReq struct {
	ID        int64  `json:"id"`
	BlackOpt  string `json:"black"`
	WhiteOpt  string `json:"white"`
	Move      string `json:"move"`
	Operation string `json:"operation"`
}

func buildBoardArr(s string) []string {
	board := make([]string, 0)

	for i := range s {
		board = append(board, string(s[i]))
	}

	return board
}
func handler(w http.ResponseWriter, r *http.Request) {

	board := buildBoardArr("...........................wb......bw...........................")

	funcMap := template.FuncMap{
		"buildBoard": func(s string) []string {
			finalBoard := make([]string, 0)
			for i := range s {
				finalBoard = append(finalBoard, fmt.Sprintf(".box data=%s\n", string(s[i])))
			}
			return finalBoard
		},
	}

	tpl, err := ace.Load("./templates/base", "./templates/board", &ace.Options{
		DynamicReload: true,
		FuncMap:       funcMap,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		// Read body
		var reqData OthelloReq
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if reqData.Operation == START_OP {
			if reqData.BlackOpt == HUMAN && reqData.WhiteOpt == HUMAN {
				humanVHuman()
			}

			if reqData.BlackOpt != HUMAN && reqData.WhiteOpt != HUMAN {
				cpuVCpu()
			}

			if reqData.BlackOpt == HUMAN && reqData.WhiteOpt != HUMAN || reqData.BlackOpt != HUMAN && reqData.WhiteOpt == HUMAN {
				humanVCpu()
			}
		}

		if reqData.Operation == MOVE_OP {

		}

		if reqData.Operation == REQUEST_OP {

		}

		if reqData.Operation == LEGAL_MOVES {

		}

	}

	if err := tpl.Execute(w, map[string][]string{"board": board}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func humanVHuman() {
	fmt.Println("...Human vs Human...")
}

func cpuVCpu() {
	fmt.Println("...AI vs AI...")
}

func humanVCpu() {
	fmt.Println("...Human vs AI...")
}

func setup(option int8) (s gomcts.GameState, p gomcts.RolloutPolicy) {

	var state gomcts.GameState = gomcts.CreateOthelloInitialGameState()
	var policy gomcts.RolloutPolicy

	if option == 0 {
		policy = gomcts.OthelloRandomRolloutPolicy
	} else if option == 1 {
		policy = gomcts.OthelloRandomRolloutPolicy
	}

	return state, policy
}

func main() {
	fsStatic := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
