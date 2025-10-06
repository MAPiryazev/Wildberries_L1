package builtins

import (
	"fmt"
	"strconv"
	"syscall"
)

// Kill отправляет сигнал SIGKILL процессу с указанным PID
func Kill(args []string) {
	if len(args) < 2 {
		fmt.Println("kill: укажите PID")
		return
	}

	pid, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("kill: ", err)
		return
	}

	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		fmt.Println("kill: ", err)
	}

}
