package storage

import (
	"errors"
	"server-calendar/internal/storage/entity"
	"sync"
	"time"
)

var (
	// ErrNotFound возвращается, если событие не найдено.
	ErrNotFound = errors.New("event not found")
	// ErrAlreadyExist возвращается, если событие с таким EventID уже существует.
	ErrAlreadyExist = errors.New("event already exists")
	// ErrUserExist возвращается, если пользователь не существует
	ErrUserExist = errors.New("user not exist")
)

type userEvent struct {
	mu     sync.RWMutex
	events map[entity.EventID]entity.Event
}

// EventStorage — простое хранилище user'ов.
type EventStorage struct {
	mu         sync.RWMutex
	userEvents map[entity.UserID]*userEvent
}

// NewStorage создаёт и возвращает новое пустое хранилище событий.
func NewStorage() *EventStorage {
	return &EventStorage{
		mu:         sync.RWMutex{},
		userEvents: make(map[entity.UserID]*userEvent),
	}
}

// CreateEvent добавляет новое событие в хранилище.
func (es *EventStorage) CreateEvent(event entity.Event) error {
	es.mu.Lock()

	_, ok := es.userEvents[event.UserID]
	if !ok {
		es.userEvents[event.UserID] = &userEvent{
			mu:     sync.RWMutex{},
			events: make(map[entity.EventID]entity.Event),
		}
	}
	es.mu.Unlock()
	es.userEvents[event.UserID].mu.Lock()
	defer es.userEvents[event.UserID].mu.Unlock()

	_, ok = es.userEvents[event.UserID].events[event.EventID]
	if !ok {
		es.userEvents[event.UserID].events[event.EventID] = event
	} else {
		return ErrAlreadyExist
	}
	return nil
}

// UpdateEvent обновляет существующее событие.
func (es *EventStorage) UpdateEvent(event entity.Event) error {
	eventsMap, ok := es.userEvents[event.UserID]
	if !ok {
		return ErrNotFound
	}
	eventsMap.mu.Lock()
	defer eventsMap.mu.Unlock()
	_, ok = eventsMap.events[event.EventID]
	if !ok {
		return ErrNotFound
	}
	eventsMap.events[event.EventID] = event
	return nil
}

// DeleteEvent удаляет событие по EventID.
func (es *EventStorage) DeleteEvent(event entity.Event) error {
	eventsMap, ok := es.userEvents[event.UserID]
	if !ok {
		return ErrNotFound
	}
	delete(eventsMap.events, event.EventID)

	return nil
}

// GetEventsByDateRange возвращает события пользователя за день.
func (es *EventStorage) GetEventsByDateRange(userID entity.UserID, from, to time.Time) ([]entity.Event, error) {
	eventsMap, ok := es.userEvents[userID]
	if !ok {
		return nil, ErrNotFound
	}
	out := eventsMap.filter(from, to)
	return out, nil
}

func (es *EventStorage) ArchiveOldEvents(cutoff time.Time) int {
	es.mu.RLock()
	defer es.mu.RUnlock()

	archived := 0

	for _, ue := range es.userEvents {
		ue.mu.Lock()
		for id, ev := range ue.events {
			if ev.Archived {
				continue
			}
			if ev.Date.Before(cutoff) {
				ev.Archived = true
				ue.events[id] = ev
				archived++
			}
		}
		ue.mu.Unlock()
	}

	return archived
}

func (u *userEvent) filter(from, to time.Time) []entity.Event {
	u.mu.RLock()
	defer u.mu.RUnlock()

	res := make([]entity.Event, 0, len(u.events)/4)
	for _, e := range u.events {
		if e.Archived {
			continue
		}
		if !e.Date.Before(from) && !e.Date.After(to) {
			res = append(res, e)
		}
	}
	return res
}
