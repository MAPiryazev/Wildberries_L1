# unixsort

Небольшая утилита на Go, которая работает похоже на стандартную команду `sort` в Linux.  
Можно сортировать строки из файла или из `stdin` с разными флагами.
Есть тесты, go vet и go lint был проверен

## Запуск

Примеры:
```bash
# обычная сортировка строк
go run main.go file.txt

# сортировка чисел
go run main.go -n file.txt

# сортировка по второй колонке (колонки разделяются табом)
go run main.go -k 2 file.txt

# уникальные строки
go run main.go -u file.txt

# реверс
go run main.go -r file.txt
