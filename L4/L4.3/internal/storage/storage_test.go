package storage_test

import (
	"errors"
	"server-calendar/internal/storage"
	"server-calendar/internal/storage/entity"
	"testing"
	"time"
)

func TestStorage_CreateEvent(t *testing.T) {
	st := storage.NewStorage()
	now := time.Now()
	str := "Test"
	ev := entity.Event{
		EventID: 1,
		UserID:  10,
		Date:    now,
		Title:   str,
	}

	err := st.CreateEvent(ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// повторное создание → ошибка
	err = st.CreateEvent(ev)
	if !errors.Is(err, storage.ErrAlreadyExist) {
		t.Fatalf("expected ErrAlreadyExist, got %v", err)
	}
}

func TestStorage_GetEvent(t *testing.T) {
	st := storage.NewStorage()
	now := time.Now().AddDate(0, 0, 1)
	ev := entity.Event{
		EventID: 1,
		UserID:  5,
		Date:    now,
	}

	_ = st.CreateEvent(ev)

	got, err := st.GetEventsByDateRange(ev.UserID, time.Now(), ev.Date.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got[0].EventID != ev.EventID {
		t.Fatalf("expected %d, got %d", ev.EventID, got[0].EventID)
	}
}

func TestStorage_GetEventsByDateRange(t *testing.T) {
	st := storage.NewStorage()

	now := time.Now()
	ev := entity.Event{
		EventID: 1,
		UserID:  77,
		Date:    now,
	}

	_ = st.CreateEvent(ev)

	from := now.Add(-time.Hour)
	to := now.Add(time.Hour)

	list, err := st.GetEventsByDateRange(ev.UserID, from, to)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 event, got %d", len(list))
	}
}
