package types

import (
	"fmt"
	"strconv"
)

type MetricFloat float64

func (m *MetricFloat) UnmarshalText(text []byte) error {
	parsed, err := parseMetricFloatString(string(text))
	if err != nil {
		return err
	}

	*m = MetricFloat(parsed)

	return nil
}

func (m MetricFloat) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%f", m), nil
}

func parseMetricFloatString(input string) (float64, error) {
	length := len(input)
	if length == 0 {
		return 0, nil
	}

	// pull off the last rune
	last := input[length-1]
	base := input[:length-1]

	var mul float64

	switch last {
	case 'k', 'K':
		mul = 1_000
	case 'm', 'M':
		mul = 1_000_000
	case 'b', 'B':
		mul = 1_000_000_000
	default:
		mul = 1
		base = input
	}

	n, err := strconv.ParseFloat(base, 64)
	if err != nil {
		return 0, fmt.Errorf("%w (base=%q): %w", ErrFailedToParseFloat, base, err)
	}

	return n * mul, nil
}
