package main

import (
	"fmt"

	"github.com/nexryai/polyxia/server/geo"
)

func main() {
	test, err := geo.LatLonToCity(35.5112036, 139.5133511)
	if err != nil {
		panic(err)
	}

	fmt.Printf("City: %s\n", test.ForecastCode)
}