package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

// Ps запускает внешнюю команду ps и печатает список процессов
func Ps(args []string) {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("ps: ошибка: ", err)
	}
}
