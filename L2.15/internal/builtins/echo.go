package builtins

import (
	"fmt"
	"strings"
)

// Echo выводит переданную строку в терминал
func Echo(args []string) {
	if len(args) > 1 {
		fmt.Println(strings.Join(args[1:], " "))
	} else {
		fmt.Println()
	}
}
