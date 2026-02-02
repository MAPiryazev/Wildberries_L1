package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "JFJFsFOIFHEIOHfhfOWFwjkfHKHfhdlHUIBIFBEUUFBPWPBFkfhdjkfhjksdhjfhjlsdhjfhjsdlfhjsdjkllfKJHFSHfL"
	mapa := make(map[rune]struct{})

	str = strings.ToLower(str)
	for _, val := range str {
		_, ok := mapa[val]
		if ok {
			fmt.Println("false")
			return
		} else {
			mapa[val] = struct{}{}
		}
	}
	fmt.Println("True")
}
