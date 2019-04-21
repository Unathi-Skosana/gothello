package main

import (
	"fmt"
	"gomcts"
	"html/template"
	"net/http"

	"github.com/yosssi/ace"
)

func buildBoardArr(s string) []string {
	board := make([]string, 0)

	for i := range s {
		board = append(board, string(s[i]))
	}

	return board
}
func handler(w http.ResponseWriter, r *http.Request) {
	initBoard := buildBoardArr("...........................wb......bw...........................")

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

	if err := tpl.Execute(w, map[string][]string{"InitBoard": initBoard}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func setup(option int8, starting int8) (s gomcts.GameState, p gomcts.RolloutPolicy) {

	var state gomcts.GameState = gomcts.CreateOthelloInitialGameState(starting)
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
