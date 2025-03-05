package details

import (
	"fmt"
	"net/http"

	"github.com/nexryai/quake/server/controller"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query().Get("id")
	if queryParams == "" {
		w.WriteHeader(400)
		return
	}

	data, err := controller.GetEventDetailsJson(queryParams)
	if err != nil || data == nil {
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", *data)
}