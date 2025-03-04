package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nexryai/quake/internal/jmafeed"
)

type Events struct {
	Events []string `json:"events"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	events, _ := jmafeed.GetJMAEvents()

	data := &Events {
		Events: *events,
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%s", json)
}