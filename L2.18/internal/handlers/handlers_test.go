package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"calendar/internal/handlers"
	"calendar/internal/models"
	"calendar/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler, _ := handlers.NewDefaultHandler(mockSvc)

	event := models.Event{
		ID:     1,
		UserID: 1,
		Date:   time.Now().Add(24 * time.Hour),
		Title:  "Встреча",
	}

	// ожидаем вызов CreateEvent и возвращаем успешный результат
	mockSvc.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(&event, nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"result"`)
}

func TestCreateEvent_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler, _ := handlers.NewDefaultHandler(mockSvc)

	// подаём неправильный JSON
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewBufferString(`{bad json`))
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid JSON")
}

func TestCreateEvent_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler, _ := handlers.NewDefaultHandler(mockSvc)

	// отсутствует title → должна быть ошибка валидации
	event := models.Event{
		UserID: 1,
		Date:   time.Now().Add(24 * time.Hour),
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateEvent(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "title cannot be empty")
}

func TestEventsForDay_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler, _ := handlers.NewDefaultHandler(mockSvc)

	userID := 1
	dateStr := "2025-10-21"

	expectedEvents := []*models.Event{
		{ID: 1, UserID: userID, Date: time.Now(), Title: "Встреча"},
	}

	mockSvc.EXPECT().
		EventsForDay(gomock.Any(), userID, dateStr).
		Return(expectedEvents, nil)

	req := httptest.NewRequest(http.MethodGet, "/events_for_day?user_id=1&date=2025-10-21", nil)
	w := httptest.NewRecorder()

	handler.EventsForDay(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"result"`)
	assert.Contains(t, w.Body.String(), `"Встреча"`)
}

func TestEventsForDay_InvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockService(ctrl)
	handler, _ := handlers.NewDefaultHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/events_for_day?user_id=1&date=invalid-date", nil)
	w := httptest.NewRecorder()

	handler.EventsForDay(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid date format")
}
