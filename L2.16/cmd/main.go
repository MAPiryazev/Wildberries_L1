package main

import (
	"flag"
	"fmt"
	"log"

	downoader "gowget/internal/downloader"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	depth := 0
	var url string

	flag.StringVar(&url, "url", "", "sets url for parsing")
	flag.IntVar(&depth, "depth", 1, "sets recursion depth for parsing")
	flag.Parse()

	if url == "" {
		fmt.Println("Укажите url через флаг -url")
		return
	}

	loader := downoader.NewStockDownloader()
	err := loader.DownloadPage(url, depth)
	if err != nil {
		fmt.Println("ошибка: ", err)
	} else {
		fmt.Println("скачивание завершено")
	}

}
