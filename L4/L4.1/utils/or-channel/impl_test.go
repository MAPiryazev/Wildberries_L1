package orchannel

import (
	"testing"
	"time"
)

func TestOrChannel_PanicsOnZeroChannels(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got nil")
		}
	}()

	// len == 0 -> log.Panicln(...) -> panic
	OrChannel()
}

func TestOrChannel_ReturnsSameChannelWhenOneProvided(t *testing.T) {
	ch := make(chan interface{})
	got := OrChannel(ch)

	// При 1 канале функция должна вернуть тот же самый канал
	if got != ch {
		t.Fatalf("expected same channel pointer, got different")
	}
}

func TestOrChannel_ClosesWhenAnyInputCloses(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	out := OrChannel(ch1, ch2)

	// Закрываем один из входных каналов -> out должен закрыться
	close(ch2)

	select {
	case _, ok := <-out:
		if ok {
			t.Fatalf("expected out to be closed")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timeout: out was not closed")
	}
}
