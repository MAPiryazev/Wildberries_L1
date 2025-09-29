package unixsort

//в файле будут лежать функции для сравнений различных типов строчных значений

func compareNumeric(a, b string, reverse bool) bool {
	na := parseNumber(a)
	nb := parseNumber(b)
	if !reverse {
		return na < nb
	}
	return na > nb
}

func compareMonth(a, b string, reverse bool) bool {
	ma := parseMonth(a)
	mb := parseMonth(b)
	if !reverse {
		return ma < mb
	}
	return ma > mb
}

func compareHumanSuff(a, b string, reverse bool) bool {
	ha := parseSuffix(a)
	hb := parseSuffix(b)
	if !reverse {
		return ha < hb
	}
	return ha > hb
}
