package main

import (
	"fmt"
	"github.com/nexryai/quake/server/controller"

	"github.com/nexryai/quake/server/jmafeed"
)

func main() {
	events, _ := jmafeed.GetJMAEvents()

	fmt.Printf("Total: %d\n", len(*events))

	for _, e := range *events {
		fmt.Println(e)
	}

	data, err := controller.GetEventDetails((*events)[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", data)
}