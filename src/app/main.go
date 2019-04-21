package main

import (
  "html/template"
  "net/http"
  "github.com/yosssi/ace"
  "othello"
)

func handler(w http.ResponseWriter, r *http.Request) {
  initBoard := " . . . . . . . . . . . . . . . . . . . . . . . . . . . w b . . . . . . b w . . . . . . . . . . . . . . . . . . . . . . . . . . . "

  funcMap := template.FuncMap{
		"buildBoard": func(s string) string {
      finalBoard := ""
      for i := range(s) {
        finalBoard += string(s[i])
      }
      return finalBoard
		},
}

  tpl, err := ace.Load("./templates/base", "./templates/inner", &ace.Options{
    DynamicReload: true,
    FuncMap: funcMap,
  })

  if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
  }


  if err := tpl.Execute(w, map[string]string{"InitBoard": initBoard}); err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
  }
}



func setup(option int8, starting int8) (s gomcts.GameState, p gomcts.RolloutPolicy) {

  var state gomcts.GameState = gomcts.CreateOthelloInitialGameState(starting)
  var policy gomcts.RolloutPolicy

  if (option == 0) {
    policy = gomcts.OthelloRandomRolloutPolicy
  } else if (option == 1) {
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
