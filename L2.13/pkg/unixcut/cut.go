package unixcut

import "strings"

// WorkLines - основная функция которая парсит строки и выводит из в stdout в соответствии с флагами
func WorkLines(lines []string, flags Flags) ([][]string, error) {
	if len(lines) == 0 {
		return [][]string{}, nil
	}

	wasEmpty := true
	if len(flags.Fields) > 0 {
		wasEmpty = false
	}

	res := make([][]string, 0)

	for _, line := range lines {
		localFields := make([]int, len(flags.Fields))
		copy(localFields, flags.Fields)

		currString := make([]string, 0)
		separatedLine := strings.Split(line, flags.Delimiter)

		if flags.Separated && len(separatedLine) <= 1 {
			continue
		}
		for i, sepVal := range separatedLine {
			if len(localFields) > 0 {
				if localFields[0] == i {
					currString = append(currString, sepVal)
					localFields = localFields[1:]
					continue
				}
			} else {
				if !wasEmpty {
					continue
				}
				currString = append(currString, sepVal)
			}
		}
		res = append(res, currString)
	}
	return res, nil

}
