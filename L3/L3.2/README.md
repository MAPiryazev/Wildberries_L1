# URL Shortener (PSQL + ClickHouse)

## Что это
Мини-сервис сокращения ссылок:
- POST /shorten — создать короткую ссылку
- GET /s/{short} — редирект на оригинал + запись клика
- GET /analytics/{short} — аналитика (total, по дням, по User-Agent)

PSQL хранит ссылки (таблица `shortcuts`). ClickHouse хранит аналитику (`shortener_analytics.click_analytics`).

## Быстрый старт
```bash
cd Wildberries_L1/L3/L3.2
docker compose up -d
```
```bash
go mod tidy
go run ./cmd/server
```
По умолчанию сервер слушает порт из `.env` (`API_PORT`, дефолт 8080).

## Эндпоинты
- POST /shorten
```json
{
  "original_url": "https://example.com",
  "client_id": "b23a8-4f12-8e11-9912-5c11c9a04f9d",
  "ttl_seconds": 3600
}
```
Ответ:
```json
{ "short_code": "abc123" }
```

- GET /s/{short}
Заголовок `X-Client-Id: <uuid>` (или query `?client_id=...`) — для аналитики. Возвращает 302 на оригинальный URL.

- GET /analytics/{short}
```json
{
  "total": 10,
  "daily": [{"key":"2025-10-29","count":3}],
  "by_user_agent": [{"key":"Mozilla/5.0","count":7}]
}
```
## Примечания
- ClickHouse init ожидает готовность через `clickhouse-client`, без curl
- Миграции для PSQL/ClickHouse применяются автоматически из docker compose


