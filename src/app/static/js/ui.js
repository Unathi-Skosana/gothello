// Swap button

export function _swap_action() {
  document.getElementById("swap-btn").onclick = evt => {
    var disabled = evt.target.attributes.disabled;

    if (disabled) {
      return;
    }

    var blackSelect = document.getElementById("black").value;
    var whiteSelect = document.getElementById("white").value;

    document.getElementById("white").value = blackSelect;
    document.getElementById("black").value = whiteSelect;
  };
}

// Start button

export function _start_action(nextToMove, actions) {
  document.getElementById("start").onclick = evt => {
    var disabled = evt.target.attributes.disabled;

    if (disabled) {
      return;
    }

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

    document.getElementById("swap-btn").setAttribute("disabled", true);
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
  };
}

// restart button
export function _restart_action(nextToMove, actions) {
  document.getElementById("restart").onclick = evt => {
    document.getElementById("white").removeAttribute("disabled");
    document.getElementById("black").removeAttribute("disabled");

    document.getElementById("swap-btn").removeAttribute("disabled");
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
  };
}
