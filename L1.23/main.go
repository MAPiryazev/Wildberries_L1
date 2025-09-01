package main

import "fmt"

func removeElement(slice []int, n int) []int {
	if n < 0 || n >= len(slice) {
		fmt.Println("Недействительный индекс удаляемого элемента")
		return slice
	}
	return append(slice[:n], slice[n+1:]...)

}

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(slice)
	slice = removeElement(slice, 3)
	fmt.Println(slice)
}
