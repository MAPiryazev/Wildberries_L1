package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// Генерируем случайное большое число
	bigRandomNum, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil))
	fmt.Printf("Случайное число: %s\n", bigRandomNum.String())

	a := big.NewInt(4_000_000_000)
	b := big.NewInt(2_000_000_000)
	fmt.Printf("a = %s\n", a.String())
	fmt.Printf("b = %s\n\n", b.String())

	sum := new(big.Int).Add(a, b)
	fmt.Printf("a + b = %s\n", sum.String())

	sub := new(big.Int).Sub(a, b)
	fmt.Printf("a - b = %s\n", sub.String())

	mult := new(big.Int).Mul(a, b)
	fmt.Printf("a * b = %s\n", mult.String())

	div := new(big.Int).Div(a, b)
	fmt.Printf("a / b = %s\n", div.String())

}
