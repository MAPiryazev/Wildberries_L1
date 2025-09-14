package main

import (
	"fmt"
	"math/rand"
	"sort"
)

const ARRAY_LENGTH = 10_000

func binarySearch(array []int, target, i, j int) int {
	if i > j {
		return -1
	}
	idx := (j + i) / 2
	if target == array[idx] {
		return idx
	} else if target > array[idx] {
		return binarySearch(array, target, idx+1, j)
	} else {
		return binarySearch(array, target, i, idx-1)
	}
}

func main() {
	array := make([]int, 0, ARRAY_LENGTH)
	for i := 0; i < ARRAY_LENGTH; i++ {
		array = append(array, rand.Int())
	}
	array = append(array, 5)

	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})

	fmt.Println(binarySearch(array, 5, 0, 10_000))

}
