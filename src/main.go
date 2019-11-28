package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/unathi-skosana/gothello/src/gomcts"
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

var (
	firstPlayer  int
	secondPlayer int
	gameState    gomcts.GameState
)

func buildBoardArr(s string) []string {
	board := make([]string, 0)

	for i := range s {
		board = append(board, string(s[i]))
	}

	return board
}

func legalActionsStr(state gomcts.GameState) []string {
	legalActions := state.GetLegalActions()
	legalActionsArr := make([]string, 0)

	for i := range legalActions {
		legalActionsArr = append(legalActionsArr, strconv.Itoa(int(legalActions[i].(gomcts.OthelloBoardGameAction).GetMove())))
	}

	return legalActionsArr

}

func handler(w http.ResponseWriter, r *http.Request) {
	var s gomcts.GameState = gomcts.CreateOthelloInitialGameState()
	board := buildBoardArr(s.(gomcts.OthelloGameState).ToString())
	legalActions := legalActionsStr(s)
	nextToMove := []string{strconv.Itoa(int(s.(gomcts.OthelloGameState).NextToMove()))}

	funcMap := template.FuncMap{
		"Combine": func(val string, id int) map[string]string {
			m := make(map[string]string)
			row := id / 8
			col := id % 8
			m["val"] = val
			m["id"] = strconv.Itoa(col + 1 + (row+1)*10)
			return m
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

	// Starting game

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

			if reqData.BlackOpt == HUMAN && reqData.WhiteOpt != HUMAN {
			}

			if reqData.BlackOpt != HUMAN && reqData.WhiteOpt == HUMAN {
				chosenAction := gomcts.MonteCarloTreeSearch(s, gomcts.OthelloHeuristicRolloutPolicy, 100)
				s = chosenAction.ApplyTo(s)
				board := buildBoardArr(s.(gomcts.OthelloGameState).ToString())
				legalActions := legalActionsStr(s)
				nextToMove := []string{strconv.Itoa(int(s.(gomcts.OthelloGameState).NextToMove()))}
			}
		}

		if reqData.Operation == MOVE_OP {

		}

		if reqData.Operation == REQUEST_OP {

		}

		if reqData.Operation == LEGAL_MOVES {

		}

	}

	// Initial state

	if err := tpl.Execute(w, map[string][]string{"board": board, "actions": legalActions, "nextToMove": nextToMove}); err != nil {
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
		policy = gomcts.OthelloHeuristicRolloutPolicy
	} else {
		policy = nil
	}

	return state, policy
}

func main() {
	fsStatic := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
