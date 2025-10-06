package main

import (
	"bufio"
	"fmt"
	"minishell/internal/core"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	shell := core.NewCore()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for range sigs {
			shell.Interrupt()
			fmt.Println()
		}
	}()

	for true {
		fmt.Print("minishell> ")
		if !reader.Scan() {
			fmt.Println("\nExit")
			break
		}

		line := reader.Text()
		shell.ExecuteLine(line)
	}
}
