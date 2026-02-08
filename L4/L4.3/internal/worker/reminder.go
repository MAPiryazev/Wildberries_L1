package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Reminder struct {
	EventID  int
	UserID   int
	Title    string
	RemindAt time.Time
}

type ReminderWorker struct {
	ch chan Reminder
}

func NewReminderWorker(buffer int) *ReminderWorker {
	return &ReminderWorker{
		ch: make(chan Reminder, buffer),
	}
}

func (w *ReminderWorker) Start() {
	go func() {
		for r := range w.ch {
			r := r // защита от capture
			delay := time.Until(r.RemindAt)
			if delay <= 0 {
				continue
			}

			time.AfterFunc(delay, func() {
				logrus.WithFields(logrus.Fields{
					"event_id": r.EventID,
					"user_id":  r.UserID,
					"title":    r.Title,
				}).Info("event reminder")
			})
		}
	}()
}

func (w *ReminderWorker) Add(r Reminder) {
	select {
	case w.ch <- r:
	default:
	}
}

func (w *ReminderWorker) Stop() {
	close(w.ch)
}
