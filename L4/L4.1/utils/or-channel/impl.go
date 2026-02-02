package orchannel

import (
	"log"
)

func OrChannel(channels ...chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		log.Panicln("не переданы каналы для обхединения")
	case 1:
		return channels[0]
	}

	out := make(chan interface{})

	go func() {
		defer close(out)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-OrChannel(append(channels[3:], out)...):
			}
		}
	}()
	return out
}
