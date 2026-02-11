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

// HumanReadable returns a human-readable string representation of the metric.
// Examples: 1000 -> "1K", 1500000 -> "1.5M", 2000000000 -> "2B".
func (m MetricFloat) HumanReadable() string {
	v := float64(m)

	if v == 0 {
		return "0"
	}

	sign := ""

	absV := v
	if v < 0 {
		absV = -v
		sign = "-"
	}

	switch {
	case absV >= 1_000_000_000:
		val := absV / 1_000_000_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dB", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fB", sign, val)
	case absV >= 1_000_000:
		val := absV / 1_000_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dM", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fM", sign, val)
	case absV >= 1_000:
		val := absV / 1_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dK", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fK", sign, val)
	default:
		if absV == float64(int64(absV)) {
			return fmt.Sprintf("%s%d", sign, int64(absV))
		}

		return fmt.Sprintf("%s%.1f", sign, absV)
	}
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
