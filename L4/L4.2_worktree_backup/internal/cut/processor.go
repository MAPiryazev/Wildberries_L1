package cut

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

func NewProcessor(delimiter string, fields []int, suppressNoDelim bool) (Processor, error) {
	if len(delimiter) == 0 {
		return nil, fmt.Errorf("%s", *ErrInvalidDelim)
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("%s", *ErrEmptyFields)
	}

	slices.Sort(fields)

	return &processor{
		cfg: CutConfig{
			Delimiter:           delimiter,
			Fields:              fields,
			SuppressNoDelimiter: suppressNoDelim,
		},
	}, nil
}

func (p *processor) ProcessLine(line string) (string, error) {
	parts := strings.Split(line, p.cfg.Delimiter)

	if len(parts) == 1 && p.cfg.SuppressNoDelimiter {
		return "", nil
	}

	var result []string
	for _, fieldIdx := range p.cfg.Fields {
		if fieldIdx < 1 {
			return "", fmt.Errorf("field index must be >= 1, got %d", fieldIdx)
		}

		idx := fieldIdx - 1
		if idx >= len(parts) {
			continue
		}
		result = append(result, parts[idx])
	}

	return strings.Join(result, p.cfg.Delimiter), nil
}

func (p *processor) ProcessReader(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		output, err := p.ProcessLine(line)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, output); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func ParseFields(fieldsStr string) ([]int, error) {
	if len(fieldsStr) == 0 {
		return nil, fmt.Errorf("%s", *ErrEmptyFields)
	}

	var fields []int
	ranges := strings.Split(fieldsStr, ",")

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("%s: invalid range %s", *ErrInvalidRange, r)
			}

			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", *ErrInvalidFields, err)
			}

			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", *ErrInvalidFields, err)
			}

			if start > end {
				return nil, fmt.Errorf("%s: start > end in range %s", *ErrInvalidRange, r)
			}

			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			field, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", *ErrInvalidFields, err)
			}
			fields = append(fields, field)
		}
	}

	slices.Sort(fields)
	fields = slices.Compact(fields)

	return fields, nil
}
