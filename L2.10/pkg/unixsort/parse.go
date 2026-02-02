package unixsort

import (
	"strconv"
	"strings"
	"unicode"
)

// функция для конвертации строки в число
func parseNumber(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

func parseMonth(s string) int {
	num, ok := monthCodes[s]
	if !ok {
		return 0
	}
	return num

}

func parseSuffix(s string) int64 {
	number := make([]byte, 0)
	suffix := make([]byte, 0)
	for i := 0; i < len(s); i++ {
		if unicode.IsDigit(rune(s[i])) {
			number = append(number, s[i])
		} else {
			suffix = append(suffix, s[i])
		}
	}
	suffix = []byte(strings.ToUpper(string(suffix)))
	intNumber, err := strconv.Atoi(string(number))
	if err != nil {
		return 0
	}
	return int64(intNumber) * sizeSuffixes[string(suffix)]
}

func trim(s []string) {
	for i, val := range s {
		s[i] = strings.TrimSpace(val)
	}
}

// функция для полученя определенной колонки в строке
func getColumn(line string, col int) string {
	if col <= 0 {
		return line
	}
	parts := strings.Split(line, "\t")
	if col > len(parts) {
		return ""
	}
	return parts[col-1]

}
