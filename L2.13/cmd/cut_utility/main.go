package main

// go run main.go -f 1,2 -d ","
// go run main.go -f 1,3-4 -d "," file.txt

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"L2.13/pkg/unixcut"
)

func main() {
	flags := unixcut.Flags{}
	var fieldsFlagString string

	flag.StringVar(&fieldsFlagString, "f", "", "номера колонок через запятую или диапазоны, например: 1,3-5")
	flag.StringVar(&flags.Delimiter, "d", "\t", "использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').")
	flag.BoolVar(&flags.Separated, "s", false, "(separated) только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).")
	flag.Parse()

	var err error
	flags.Fields, err = unixcut.ParseFields(fieldsFlagString)
	if err != nil {
		log.Panicln("неверный формат ввода", err)
	}

	lines := make([]string, 0)
	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		file, err := os.Open(filename)
		if err != nil {
			log.Panicln(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))

		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))
		}
	}

	out, err := unixcut.WorkLines(lines, flags)

	for _, val := range out {
		fmt.Println(val)
	}

}
