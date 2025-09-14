package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	temperatures := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	mapa := make(map[int][]float64)
	sort.Slice(temperatures, func(i, j int) bool {
		return temperatures[i] < temperatures[j]
	})

	for _, val := range temperatures {
		group := math.Floor(val/10) * 10
		mapa[int(group)] = append(mapa[int(group)], val)
	}

	fmt.Println(mapa)

}
