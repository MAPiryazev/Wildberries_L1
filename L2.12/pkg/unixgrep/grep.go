package unixgrep

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func HandleGrep(lines []string, pattern string, flags Flags) []string {
	searchMode := func(s string) bool {
		if flags.IgnoreCase { // -i
			s = strings.ToLower(s)
			pattern = strings.ToLower(pattern)
		}
		if flags.FixedString { // -F
			return strings.Contains(s, pattern)
		} else {
			re := regexp.MustCompile(pattern)
			return re.MatchString(s)
		}
	}

	matches := findMatches(lines, flags, searchMode)
	return getMatches(lines, flags, matches)

}

func findMatches(lines []string, flags Flags, searchMode func(s string) bool) []bool {
	matches := make([]bool, len(lines))
	for i, str := range lines {
		matches[i] = searchMode(str)
	}
	if flags.InvertMatch { // -v
		for i := 0; i < len(matches); i++ {
			matches[i] = !matches[i]
		}
	}
	return matches
}

func getMatches(lines []string, flags Flags, matches []bool) []string {
	if flags.CountOnly { // -c
		count := 0
		for _, val := range matches {
			if val {
				count++
			}
		}
		return []string{strconv.Itoa(count)}
	}

	include := make(map[int]bool)
	for i, ok := range matches {
		if ok {
			include[i] = true
		}

		// -B
		for j := i - flags.Before; j < i; j++ {
			if j >= 0 {
				include[j] = true
			}
		}
		// -A
		for j := i + 1; j <= i+flags.After && j < len(matches); j++ {
			include[j] = true
		}
		//- C
		for j := i - flags.Context; j < i+flags.Context; j++ {
			if j >= 0 && j < len(matches) {
				include[j] = true
			}
		}
	}

	var result []string
	for i := 0; i < len(lines); i++ {
		if include[i] {
			if flags.LineNumer { // -n
				result = append(result, fmt.Sprintf("%v:%v", i+1, lines[i]))
			} else {
				result = append(result, lines[i])
			}
		}
	}

	return result

}
