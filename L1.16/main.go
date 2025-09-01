package main

import (
	"fmt"
	"math/rand"
)

func quicksort(array []int, left, right int) {
	if left >= right {
		return //в этом случае массив уже отсортирован
	}
	pivot := array[(left+right)/2]
	i, j := left, right
	for i <= j {
		for array[i] < pivot && i <= right {
			i++
		}
		for array[j] > pivot && i >= left {
			j--
		}
		if i <= j {
			array[i], array[j] = array[j], array[i]
			i++
			j--
		}
	}
	quicksort(array, left, j)
	quicksort(array, i, right)
}

func main() {
	array := make([]int, 0, 10_000)
	for i := 0; i < 10_000; i++ {
		array = append(array, rand.Int())
	}
	quicksort(array, 0, len(array)-1)
	fmt.Println(array)
}
