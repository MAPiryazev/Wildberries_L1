package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func readString() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	str := ""
	if scanner.Scan() {
		str = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return str, nil
}

func writeSequence(sb *strings.Builder, val rune, amount rune) {
	count := int(amount - '0')
	for i := 0; i < count-1; i++ {
		sb.WriteRune(val)
	}
}

func unpackSrting(input string) (string, error) {
	if input == "" {
		return input, nil
	}
	var sb strings.Builder
	sb.Grow(len(input))

	escape := false
	var lastChar rune
	for _, val := range input {

		if val == rune('\\') {
			escape = true
			continue
		}

		if escape {
			sb.WriteRune(val)
			lastChar = val
			escape = false
			continue
		}

		if unicode.IsDigit(val) && !escape {
			if lastChar == 0 {
				return "", fmt.Errorf("неправильный формат введенной строки")
			}
			writeSequence(&sb, lastChar, val)
			continue
		}

		if !unicode.IsDigit(val) {
			sb.WriteRune(val)
			lastChar = val
			continue
		}
	}
	return sb.String(), nil

}

func main() {
	input, err := readString()
	if err != nil {
		log.Panicln("ошибка при считывании строки ", err)
	}
	unpackedString, err := unpackSrting(input)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(unpackedString)

}
