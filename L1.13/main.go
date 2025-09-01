package main

import "fmt"

func main() {
	var a, b int
	fmt.Println("Введите первое число")
	fmt.Scan(&a)
	fmt.Println("Введите второе число")
	fmt.Scan(&b)

	a = a ^ b
	b = a ^ b
	a = a ^ b
	fmt.Printf("a: %d, b: %d", a, b)
}
