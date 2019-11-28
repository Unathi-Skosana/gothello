var HUMAN = 1;
var RANDOM_AI = 2;
var SMART_AI = 3;

document.addEventListener("DOMContentLoaded", function() {
  console.log("DOCUMENT LOADED");
});

// Swap button

function _swap_action() {
  var disabled = document.getElementById("swap").attributes.disabled;

  if (disabled) return;

  var blackSelect = document.getElementById("black").value;
  var whiteSelect = document.getElementById("white").value;

  document.getElementById("white").value = blackSelect;
  document.getElementById("black").value = whiteSelect;
}

// Start button
function _start_action(nextToMove, actions) {
  var disabled = document.getElementById("start").attributes.disabled;

  if (disabled) return;

  for (key in actions) {
    if (nextToMove == "1") {
      document
        .getElementById(actions[key])
        .classList.add("black-possible-move");
    }
    if (nextToMove == "-1") {
      document
        .getElementById(actions[key])
        .classList.add("white-possible-move");
    }
  }

  var blackSelect = document.getElementById("black").value;
  var whiteSelect = document.getElementById("white").value;

  document.getElementById("white").setAttribute("disabled", true);
  document.getElementById("black").setAttribute("disabled", true);

  document.getElementById("swap").setAttribute("disabled", true);
  document.getElementById("start").setAttribute("disabled", true);

  var xhr = new XMLHttpRequest();

  xhr.open("POST", "http://localhost:8080", true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.send(
    JSON.stringify({
      operation: "start",
      black: blackSelect,
      white: whiteSelect,
      move: ""
    })
  );
}

// restart button
function _restart_action(nextToMove, actions) {
  document.getElementById("white").removeAttribute("disabled");
  document.getElementById("black").removeAttribute("disabled");

  document.getElementById("swap").removeAttribute("disabled");
  document.getElementById("start").removeAttribute("disabled");

  for (key in actions) {
    if (nextToMove == "1") {
      document
        .getElementById(actions[key])
        .classList.remove("black-possible-move");
    }
    if (nextToMove == "-1") {
      document
        .getElementById(action[key])
        .classList.remove("white-possible-move");
    }
  }

  var xhr = new XMLHttpRequest();
  xhr.open("POST", "http://localhost:8080", true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.send(
    JSON.stringify({
      operation: "restart",
      black: "",
      white: "",
      move: ""
    })
  );
}
