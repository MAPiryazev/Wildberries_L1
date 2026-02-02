package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func run(address string, timeout time.Duration) {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprintf(os.Stderr, "Ошибка: не удалось подключиться к %s за %s\n", address, timeout.String())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Ошибка подключения к %s: %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()

	var tcpConn *net.TCPConn
	if c, ok := conn.(*net.TCPConn); ok {
		tcpConn = c
	}

	done := make(chan struct{})
	var once sync.Once
	closeDone := func() { once.Do(func() { close(done) }) }

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigs:
			_ = conn.Close()
			closeDone()
		case <-done:
			return
		}
	}()

	go func() {
		defer closeDone()
		writer := bufio.NewWriter(os.Stdout)
		defer writer.Flush()
		_, _ = io.Copy(writer, conn)
	}()

	reader := bufio.NewReader(os.Stdin)
	if timeout > 0 {
		_ = conn.SetDeadline(time.Now().Add(timeout))
	}
	_, err = io.Copy(conn, reader)
	if err == nil || errors.Is(err, io.EOF) {
		if tcpConn != nil {
			_ = tcpConn.CloseWrite()
		} else {
			_ = conn.Close()
		}
	}

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}
}
