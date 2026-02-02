package main

// go run main.go -i hello file.txt
// go run main.go -n two

import (
	"flag"
	"fmt"
	"log"

	"L2.12/pkg/unixgrep"
)

func main() {
	flags := unixgrep.Flags{}

	flag.IntVar(&flags.After, "A", 0, "выводить N строк после совпадения")
	flag.IntVar(&flags.Before, "B", 0, "выводить N срок перед совпадением")
	flag.IntVar(&flags.Context, "C", 0, "выводить N строк вокруг совпадения")
	flag.BoolVar(&flags.CountOnly, "c", false, "выводить только количество совпадений")
	flag.BoolVar(&flags.IgnoreCase, "i", false, "не учитывать регистр")
	flag.BoolVar(&flags.InvertMatch, "v", false, "инвертировать фильтр")
	flag.BoolVar(&flags.FixedString, "F", false, "точное совпадение подстроки")
	flag.BoolVar(&flags.LineNumer, "n", false, "выводить номер строки")
	flag.Parse()
	flags.Normalize()

	lines := make([]string, 0)
	if flag.NArg() > 1 {
		filename := flag.Arg(1)
		var err error
		lines, err = unixgrep.ReadLinesFromFile(filename)
		if err != nil {
			log.Panicln("ошибка чтения файла", err)
		}
	} else {
		lines = unixgrep.ReadLinesFromStdin()
	}

	res := unixgrep.HandleGrep(lines, flag.Arg(0), flags)
	fmt.Println(res)

}
