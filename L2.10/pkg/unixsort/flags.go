package unixsort

type Flags struct {
	Column      int  // -k N
	Numeric     bool // -n имеется ввиду просто строка из чисел
	Reverse     bool // -r
	Unique      bool // -u
	Month       bool // -m
	IgnoreSpace bool // -b
	Check       bool // -c
	Human       bool // -h для суффиксов рзамеров типо KB, MB и т д
}
