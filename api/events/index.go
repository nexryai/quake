package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nexryai/quake/server/controller"
)

type Events struct {
	Events []string `json:"events"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	events, _ := controller.GetJMAEvents()

	data := &Events {
		Events: *events,
	}

	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", res)
}