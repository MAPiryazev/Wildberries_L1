package unixsort

import "sort"

//SortLines служит главным методом для сортировки
func SortLines(lines []string, opts Flags) []string {
	result := make([]string, len(lines))
	copy(result, lines)

	//режем пробелы если есть флаг -b
	if opts.IgnoreSpace {
		trim(result)
	}
	if opts.Check {
		checkSorted(lines, opts)
		return result
	}

	sort.Slice(result, func(i, j int) bool {
		a := getColumn(result[i], opts.Column)
		b := getColumn(result[j], opts.Column)

		if opts.Numeric { // если флаг -n (сортировка чисел)
			return compareNumeric(a, b, opts.Reverse)
		}
		if opts.Month { // если флаг -m (сортировка месяцев)
			return compareMonth(a, b, opts.Reverse)
		}
		if opts.Human { // если флаг -h сравнение с учетом сокращений
			return compareHumanSuff(a, b, opts.Reverse)
		}

		if opts.Reverse { // обычная лексикографическая сортировка
			return b < a
		}
		return a < b
	})

	if opts.Unique {
		result = makeUnique(result) // если флаг -u оставляем уникальные строки
	}

	return result

}
