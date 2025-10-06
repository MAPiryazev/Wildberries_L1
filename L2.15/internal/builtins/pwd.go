package builtins

import (
	"fmt"
	"os"
)

// Pwd печатает текущую рабочую директорию
func Pwd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("ошибка pwd: ", err)
		return
	}
	fmt.Println(dir)
}
