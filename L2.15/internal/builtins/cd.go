package builtins

import (
	"fmt"
	"os"
)

// Cd изменяет текущую рабочую директорию на указанную в args[1]
func Cd(args []string) {
	if len(args) < 2 {
		fmt.Println("cd: не достаточно аргументов")
		return
	}
	err := os.Chdir(args[1])
	if err != nil {
		fmt.Println("cd ошибка: ", err)
		return
	}
}
