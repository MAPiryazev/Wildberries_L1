package cut

import "io"

type CutConfig struct {
	Delimiter           string
	Fields              []int
	SuppressNoDelimiter bool
}

type Processor interface {
	ProcessLine(line string) (string, error)
	ProcessReader(r io.Reader, w io.Writer) error
}

type processor struct {
	cfg CutConfig
}

var (
	ErrInvalidFields  = new(string)
	ErrEmptyFields    = new(string)
	ErrInvalidDelim   = new(string)
	ErrInvalidRange   = new(string)
	ErrProcessingLine = new(string)
)

func init() {
	*ErrInvalidFields = "invalid fields format"
	*ErrEmptyFields = "fields list is empty"
	*ErrInvalidDelim = "delimiter is empty"
	*ErrInvalidRange = "field index out of range"
	*ErrProcessingLine = "error processing line"
}
