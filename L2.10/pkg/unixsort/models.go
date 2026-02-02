package unixsort

var monthCodes = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}

var sizeSuffixes = map[string]int64{
	"B":  1,
	"KB": 1 << 10,
	"MB": 1 << 20,
	"GB": 1 << 30,
	"TB": 1 << 40,
}
