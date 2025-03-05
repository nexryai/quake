package main

import (
	"fmt"

	"github.com/nexryai/quake/server/controller"
)

func main() {
	events, _ := controller.GetJMAEvents()

	fmt.Printf("Total: %d\n", len(*events))

	for _, e := range *events {
		fmt.Println(e)
	}

	data, err := controller.GetEventDetailsJson((*events)[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", *data)
}