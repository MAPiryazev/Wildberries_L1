package unixcut

// Flags структура, предназначенная для хранения флагов и передачи их в фукнции
type Flags struct {
	Fields    []int  // -f
	Delimiter string // -d
	Separated bool   // -s
}
