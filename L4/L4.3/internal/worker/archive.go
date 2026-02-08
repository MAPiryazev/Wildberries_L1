package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

type ArchiveWorker struct {
	interval time.Duration
	done     chan struct{}
	archive  func(cutoff time.Time) int
}

func NewArchiveWorker(interval time.Duration, archiveFunc func(time.Time) int) *ArchiveWorker {
	return &ArchiveWorker{
		interval: interval,
		done:     make(chan struct{}),
		archive:  archiveFunc,
	}
}

func (w *ArchiveWorker) Start() {
	ticker := time.NewTicker(w.interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cutoff := time.Now().AddDate(0, 0, -30)
				count := w.archive(cutoff)
				logrus.WithFields(logrus.Fields{
					"archived": count,
					"cutoff":   cutoff.Format(time.DateOnly),
				}).Info("archive completed")

			case <-w.done:
				return
			}
		}
	}()
}

func (w *ArchiveWorker) Stop() {
	close(w.done)
}
