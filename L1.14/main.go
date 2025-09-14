package main

import (
	"fmt"
	"reflect"
)

func displayType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println("int", v)
	case string:
		fmt.Println("string", v)
	case bool:
		fmt.Println("bool", v)
	default:
		t := reflect.TypeOf(i)
		if t.Kind() == reflect.Chan {
			fmt.Println("Канал")
		} else {
			fmt.Println("Неизвестный тип")
		}
	}
}

func main() {
	a := 0
	b := "str"
	c := false
	d := make(chan int)

	displayType(a)
	displayType(b)
	displayType(c)
	displayType(d)

}
