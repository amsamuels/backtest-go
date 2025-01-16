package data

import (
	"fmt"
	"time"
)

// Timeframe constants.
const (
	OneHour  Timeframe = "1h"
	FourHour Timeframe = "4h"
	Daily    Timeframe = "1d"
)

// ParseTimeframe converts a string into a Timeframe.
func ParseTimeframe(tf string) (Timeframe, error) {
	switch tf {
	case "1m", "1h", "4h", "1d":
		return Timeframe(tf), nil
	default:
		return "", fmt.Errorf("invalid timeframe: %s", tf)
	}
}

// Duration converts a Timeframe into a time.Duration.
func (tf Timeframe) Duration() (time.Duration, error) {
	switch tf {
	case OneHour:
		return time.Hour, nil
	case FourHour:
		return 4 * time.Hour, nil
	case Daily:
		return 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid timeframe: %s", tf)
	}
}
