package data

import "time"

// Trade represents an executed trade.
type Trade struct {
	Timestamp time.Time
	Type      string
	Price     float64
	Volume    float64
}
