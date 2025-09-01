package main

import (
	"fmt"
)

func main() {
	inputString := "🔥🐟Главрыба😊"
	runeString := []rune(inputString)

	i, j := 0, len(runeString)-1
	for i < j {
		runeString[i], runeString[j] = runeString[j], runeString[i]
		i++
		j--
	}
	fmt.Println(string(runeString))

}
