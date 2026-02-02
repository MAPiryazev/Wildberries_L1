package unixcut

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ParseFields функция, предназначенная для корректного считывания флага f и его аргументов
func ParseFields(flagString string) ([]int, error) {
	flagString = strings.TrimSpace(flagString)
	if len(flagString) == 0 {
		return []int{}, nil
	}

	res := make([]int, 0) // слайс для номеров колонок которые надо выводить
	separatedString := strings.Split(flagString, ",")
	for _, val := range separatedString {
		if strings.Contains(val, "-") {
			borders := strings.Split(val, "-")

			if len(borders) != 2 {
				return nil, fmt.Errorf("неверно заданы границы столбцов для вывода")
			}

			left, err := strconv.Atoi(borders[0])
			right, err := strconv.Atoi(borders[1])
			if err != nil {
				return nil, err
			}

			if left > right {
				return nil, fmt.Errorf("правая граница не может быть меньше чем левая")
			}

			for i := left; i <= right; i++ {
				res = append(res, i-1)
			}
		} else {
			idx, _ := strconv.Atoi(val)
			res = append(res, idx-1)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res, nil

}
