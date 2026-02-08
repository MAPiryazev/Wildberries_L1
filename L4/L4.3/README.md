# Календарь

Мини-сервис календаря: CRUD событий + выборка на день/неделю/месяц, фоновые воркеры (архивация/напоминания).

## Запуск локально
```bash
go run ./cmd/app
```

## Docker
```bash
docker build -t calendar-server .
docker run --rm -p 8081:8081 calendar-server
```