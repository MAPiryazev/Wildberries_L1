package main

import (
	"fmt"

	points "L1.24/models"
)

func main() {
	a := points.NewPoint(2, 2)
	b := points.NewPoint(4, 4)

	fmt.Println(a.Distance(b))
}
