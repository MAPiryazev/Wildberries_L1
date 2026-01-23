# L4.1
реализация паттерна **or-channel** для объединения stop-сигналов из нескольких каналов и пример использования в worker pool.

## Пример использования (main)

В `main.go`:
- `inputGen(ctx)` генерирует числа в `chan int` с тикером, пока контекст не отменён.
- `chanArrayGen(n, seconds)` создаёт `n` каналов `chan interface{}`, через `seconds` секунд отправляет значение в случайный канал и (по `defer`) закрывает каналы.
- `newWorkerPool(10, input, orchannel.OrChannel(stopChannelsArray...))` создаёт пул, который прекращает работу, когда `stopChannel` сработает.

Запуск воркеров:
- Каждый воркер в цикле делает `select`: либо `<-stopChannel` (остановка), либо читает из `inputChannel` и прибавляет к сумме под mutex.

## Запуск

Тесты:
```bash
go test ./...
