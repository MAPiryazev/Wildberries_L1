package core

import (
	"fmt"
	"minishell/internal/builtins"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// Core используется как хранилище для команд которые надо исполнить
type Core struct {
	mu         sync.Mutex
	currentCmd *exec.Cmd
}

// NewCore - конструктор для Core
func NewCore() *Core {
	return &Core{}
}

func (c *Core) setCurrentCmd(cmd *exec.Cmd) {
	c.mu.Lock()
	c.currentCmd = cmd
	c.mu.Unlock()
}

func (c *Core) clearCurrentCmd() {
	c.mu.Lock()
	c.currentCmd = nil
	c.mu.Unlock()
}

// Interrupt прервывает текущий процесс или группу
func (c *Core) Interrupt() {
	c.mu.Lock()
	cmd := c.currentCmd
	c.mu.Unlock()
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
}

// ExecuteLine парсит строку, и выполняет команды, содержащиеся в ней, поддерживает |
func (c *Core) ExecuteLine(line string) {
	if strings.Contains(line, "|") {
		segments := strings.Split(line, "|")
		var stages [][]string
		for _, seg := range segments {
			seg = strings.TrimSpace(seg)
			if seg == "" {
				continue
			}
			parts := strings.Fields(seg)
			if len(parts) == 0 {
				continue
			}
			stages = append(stages, parts)
		}
		if len(stages) == 0 {
			return
		}
		c.runPipeline(stages)
		return
	}

	args := strings.Fields(line)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "cd":
		builtins.Cd(args)
	case "pwd":
		builtins.Pwd(args)
	case "echo":
		builtins.Echo(args)
	case "kill":
		builtins.Kill(args)
	case "ps":
		builtins.Ps(args)
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		c.setCurrentCmd(cmd)
		err := cmd.Run()
		c.clearCurrentCmd()
		if err != nil {
			fmt.Println("ошибка выполнения команды: ", err)
		}
	}

}

// runPipeline запускает конвейер команд с общей группой процессов
func (c *Core) runPipeline(stages [][]string) {
	n := len(stages)
	if n == 1 {
		args := stages[0]
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		c.setCurrentCmd(cmd)
		err := cmd.Run()
		c.clearCurrentCmd()
		if err != nil {
			fmt.Println("ошибка выполнения команды: ", err)
		}
		return
	}

	cmds := make([]*exec.Cmd, 0, n)
	for _, args := range stages {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		cmds = append(cmds, cmd)
	}

	var prevRead *os.File
	for i, cmd := range cmds {
		if i == 0 {
			cmd.Stdin = os.Stdin
		} else {
			cmd.Stdin = prevRead
		}

		if i == n-1 {
			cmd.Stdout = os.Stdout
		} else {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Println("ошибка создания pipe: ", err)
				return
			}
			cmd.Stdout = w
			prevRead = r
		}
	}

	first := cmds[0]
	first.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.setCurrentCmd(first)
	if err := first.Start(); err != nil {
		c.clearCurrentCmd()
		fmt.Println("ошибка запуска команды: ", err)
		return
	}
	if w, ok := first.Stdout.(*os.File); ok {
		_ = w.Close()
	}
	var toClose []*os.File
	if r, ok := cmds[1].Stdin.(*os.File); ok && r != os.Stdin {
		toClose = append(toClose, r)
	}

	pgid := first.Process.Pid

	for i := 1; i < n; i++ {
		cmd := cmds[i]
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: pgid}
		if err := cmd.Start(); err != nil {
			fmt.Println("ошибка запуска команды: ", err)
			break
		}
		if w, ok := cmd.Stdout.(*os.File); ok {
			_ = w.Close()
		}
		if i+1 < n {
			if r, ok := cmds[i+1].Stdin.(*os.File); ok && r != os.Stdin {
				toClose = append(toClose, r)
			}
		}
	}

	for _, f := range toClose {
		_ = f.Close()
	}

	for _, cmd := range cmds {
		_ = cmd.Wait()
	}

	c.clearCurrentCmd()
}
