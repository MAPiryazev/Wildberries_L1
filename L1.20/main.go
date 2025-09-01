package main

import "fmt"

func reverseStr(str []rune) {
	left, right := 0, len(str)-1
	for left < right {
		str[left], str[right] = str[right], str[left]
		left++
		right--
	}
}

func main() {
	str := "snow dog sun"
	runeStr := []rune(str)

	reverseStr(runeStr)

	i := 0
	for j := 0; j < len(runeStr)+1; j++ {
		if j == len(runeStr) || runeStr[j] == ' ' {
			reverseStr(runeStr[i:j])
			i = j + 1
		}
	}

	fmt.Println(string(runeStr))
}
