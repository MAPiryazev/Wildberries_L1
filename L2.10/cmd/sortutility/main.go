package main

// go run main.go file.txt
// go run main.go -n file.txt
// go run main.go -k 2 file.txt сортировка по второй колонке

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"L2.10/pkg/unixsort"
)

func main() {
	column := flag.Int("k", 0, "сортировать по номеру колонки (нумерация с 1)")
	if *column < 0 {
		log.Panicln("номер колонки не может быть < 0")
	}
	numeric := flag.Bool("n", false, "сортировать числа")
	reverse := flag.Bool("r", false, "обратный порядок сортировки")
	unique := flag.Bool("u", false, "оставить только уникальные строки")
	month := flag.Bool("m", false, "сортировать по названию месяца")
	ignoreSpace := flag.Bool("b", false, "не учитывать пробелы сначала и в конце строк")
	check := flag.Bool("c", false, "проверить отсортированы ли данные")
	human := flag.Bool("h", false, "сортировка с учетом суффиксов типо MB, KB и т д")

	flag.Parse()
	opts := unixsort.Flags{
		Column:      *column,
		Numeric:     *numeric,
		Reverse:     *reverse,
		Unique:      *unique,
		Month:       *month,
		IgnoreSpace: *ignoreSpace,
		Check:       *check,
		Human:       *human,
	}

	var input *os.File
	var err error

	if len(flag.Args()) > 0 {
		input, err = os.Open(flag.Args()[0])
		if err != nil {
			log.Panicln(err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	var lines []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Panicln(err)
	}

	result := unixsort.SortLines(lines, opts)
	for _, line := range result {
		fmt.Println(line)
	}

}
