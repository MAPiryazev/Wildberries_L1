package log

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type AsyncLogger struct {
	ch chan LogEntry
}

type LogEntry struct {
	Method   string
	Path     string
	Duration time.Duration
	Time     time.Time
}

func NewAsyncLogger(buffer int) *AsyncLogger {
	return &AsyncLogger{
		ch: make(chan LogEntry, buffer),
	}
}

func (l *AsyncLogger) Start() {
	go func() {
		for entry := range l.ch {
			logrus.WithFields(logrus.Fields{
				"method":   entry.Method,
				"path":     entry.Path,
				"duration": entry.Duration,
			}).Info("http request")
		}
	}()
}

func (l *AsyncLogger) Log(entry LogEntry) {
	select {
	case l.ch <- entry:
	default:
	}
}

func (l *AsyncLogger) Stop() {
	close(l.ch)
}

func LoggerMiddleware(logger *AsyncLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			logger.Log(LogEntry{
				Method:   r.Method,
				Path:     r.URL.Path,
				Duration: time.Since(start),
				Time:     time.Now(),
			})
		})
	}
}
