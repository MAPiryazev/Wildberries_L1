package unixsort

import "fmt"

func makeUnique(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	var result []string

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func checkSorted(lines []string, opts Flags) {
	for i := 0; i < len(lines)-1; i++ {
		a := getColumn(lines[i], opts.Column)
		b := getColumn(lines[i+1], opts.Column)

		sorted := true
		if opts.Numeric {
			sorted = compareNumeric(a, b, false)
		} else if opts.Month {
			sorted = compareMonth(a, b, false)
		} else if opts.Human {
			sorted = compareHumanSuff(a, b, false)
		} else {
			sorted = a <= b
		}

		if !sorted {
			fmt.Println("Строки не отсортированы")
			return
		}
	}
	fmt.Println("Строки отсортированы")
}
