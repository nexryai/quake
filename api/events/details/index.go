package details

import (
	"fmt"
	"net/http"

	"github.com/nexryai/polyxia/server/controller"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	debug := false
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryParams := r.URL.Query().Get("id")
	if queryParams == "" {
		w.WriteHeader(400)
		return
	}

	debugFlag := r.URL.Query().Get("debug")
	if debugFlag == "dummy" {
		debug = true
	}

	data, err := controller.GetEventDetailsJson(queryParams, debug)
	if err != nil || data == nil {
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", *data)
}