package unixgrep

type Flags struct {
	After       int  // -A
	Before      int  // -B
	Context     int  // -C
	CountOnly   bool // -c
	IgnoreCase  bool // -i
	InvertMatch bool // -v
	FixedString bool // -F
	LineNumer   bool // -n
}

func (f *Flags) Normalize() {
	if f.Context > 0 {
		f.Before = 0
		f.After = 0
	}
}
