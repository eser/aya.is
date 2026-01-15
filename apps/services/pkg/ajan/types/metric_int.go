package types

import (
	"fmt"
	"math"
	"strconv"
)

type MetricInt int64

func (m *MetricInt) UnmarshalText(text []byte) error {
	parsed, err := parseMetricIntString(string(text))
	if err != nil {
		return err
	}

	*m = MetricInt(parsed)

	return nil
}

func (m MetricInt) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%d", m), nil
}

func parseMetricIntString(input string) (int64, error) {
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

	n, err := strconv.ParseFloat(base, 64) //nolint:varnamelen
	if err != nil {
		return 0, fmt.Errorf("%w (base=%q): %w", ErrFailedToParseFloat, base, err)
	}

	// FIXME(@eser) this is a hack to round the number to the nearest integer
	return int64(math.Round(n * mul)), nil
}
