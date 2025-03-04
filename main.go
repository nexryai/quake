package main

import (
	"fmt"

	"github.com/nexryai/quake/server/jmafeed"
)

func main() {
	events, _ := jmafeed.GetJMAEvents()

	fmt.Printf("Total: %d\n", len(*events))

	for _, e := range *events {
		fmt.Println(e)
	}
}