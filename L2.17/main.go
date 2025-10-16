package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Example: telnet --timeout=5s smtp.gmail.com 25")
	flag.PrintDefaults()
}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		usage()
		os.Exit(2)
	}

	host := strings.TrimSpace(args[0])
	port := strings.TrimSpace(args[1])
	address := net.JoinHostPort(host, port)

	run(address, *timeout)
}
