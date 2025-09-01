// var justString string

// func someFunc() {
//   v := createHugeString(1<<10)
//   justString = v[:100]
// }

// func main() {
//   someFunc()
// }

/*
1) утечка памяти - из-за того, что juststring ссылается только на небольшую часть v, большая часть v остается неиспользованной,
а сборщик мусора не может освободить память потому что на v есть ссылка
2) глобальная перменная juststring, если эту функцию вызывать в разных горутинах то будет гонка данных по области juststring
*/

package main

import (
	"sync"
)

type justString struct { //структура подобно потокобезопасной мапе, с мьютексом
	str string
	mu  sync.Mutex
}

func NewJustString() *justString {
	return &justString{
		str: string(make([]byte, 0, 0)),
		mu:  sync.Mutex{},
	}
}

func someFunc(j *justString) {
	v := createHugeString(1 << 10)

	temp := append(make([]byte, 0, 100), v[:100]...) //здесь теперь создается копия а не срез от v, чтобы он могу быть высвобожден GC
	j.mu.Lock()
	defer j.mu.Unlock()
	j.str = string(temp)

}

func main() {
	js := NewJustString()
	someFunc(js)
}
