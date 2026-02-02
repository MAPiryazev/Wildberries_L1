package main

import (
	"errors"
	"fmt"
	"log"
)

func writeBit(num int64, i, value int) (res int64, err error) {
	if i > 64 || i < 1 {
		return 0, errors.New("Недопустимый индекс бита для замены")
	}
	i--

	bitmask := int64(1) << i
	if value == 1 {
		num = num | bitmask
	} else {
		num = num & (^bitmask)
	}
	return num, nil
}

func main() {
	fmt.Println("Введите число")
	var num int64
	fmt.Scan(&num)
	fmt.Println("Введите номер бита")
	var i int
	fmt.Scan(&i)
	fmt.Println("Введите значение бита")
	var value int
	fmt.Scan(&value)
	if value != 0 && value != 1 {
		log.Fatal("Бит может принимать значения только 0 или 1")
	}

	res, err := writeBit(num, i, value)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Результат: ", res)

}
