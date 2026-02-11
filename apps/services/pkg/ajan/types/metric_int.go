package types

import (
	"encoding/json"
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

func (m *MetricInt) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a number
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*m = MetricInt(num)

		return nil
	}

	// If not a number, try as a string (e.g., "3400K", "1M")
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("MetricInt must be a number or string: %w", err)
	}

	parsed, err := parseMetricIntString(str)
	if err != nil {
		return err
	}

	*m = MetricInt(parsed)

	return nil
}

func (m MetricInt) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%d", m), nil
}

// HumanReadable returns a human-readable string representation of the metric.
// Examples: 1000 -> "1K", 1500000 -> "1.5M", 2000000000 -> "2B".
func (m MetricInt) HumanReadable() string {
	v := int64(m)

	if v == 0 {
		return "0"
	}

	absV := v
	sign := ""

	if v < 0 {
		absV = -v
		sign = "-"
	}

	switch {
	case absV >= 1_000_000_000:
		val := float64(absV) / 1_000_000_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dB", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fB", sign, val)
	case absV >= 1_000_000:
		val := float64(absV) / 1_000_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dM", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fM", sign, val)
	case absV >= 1_000:
		val := float64(absV) / 1_000
		if val == float64(int64(val)) {
			return fmt.Sprintf("%s%dK", sign, int64(val))
		}

		return fmt.Sprintf("%s%.1fK", sign, val)
	default:
		return fmt.Sprintf("%s%d", sign, absV)
	}
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
