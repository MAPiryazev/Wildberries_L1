package main

//сложность - O(n * m * log(m))

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type node struct {
	baseWord  string
	anagramms []string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	inputStr := strings.Fields(line)
	if len(inputStr) == 0 {
		return
	}

	mapa := make(map[string]*node)
	for _, val := range inputStr { //цикл по всем словам в строке O(n)
		sliceRune := []rune(val)
		sort.Slice(sliceRune, func(i, j int) bool { // сортировка слайсов в цикле
			// за O(m*log(m)) т.к. в go используется quickSort
			return sliceRune[i] < sliceRune[j]
		})
		_, ok := mapa[string(sliceRune)]
		if ok {
			mapa[string(sliceRune)].anagramms = append(mapa[string(sliceRune)].anagramms, val)
		} else {
			nodeObj := &node{baseWord: val, anagramms: nil}
			mapa[string(sliceRune)] = nodeObj
		}
	}

	for _, val := range mapa {
		if len(val.anagramms) == 0 {
			continue
		}
		fmt.Printf("-\"%v\":[ %v ]\n", val.baseWord, val.anagramms)
	}

}
