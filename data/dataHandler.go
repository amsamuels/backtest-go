package data

import (
	"fmt"
	"sort"
	"time"
)

// Timeframe defines the duration of a candlestick (e.g., 1m, 5m, 1h).
type Timeframe string

// DataHandler manages candlesticks grouped by timeframe.
type DataHandler struct {
	CandlesByTimeframe map[Timeframe]map[time.Time]Candle
}

// NewDataHandler initializes a DataHandler with no candles.
func NewDataHandler() *DataHandler {
	return &DataHandler{
		CandlesByTimeframe: make(map[Timeframe]map[time.Time]Candle),
	}
}

// AddCandles adds a slice of candles to a specific timeframe.
func (dh *DataHandler) AddCandles(timeframe Timeframe, candles []Candle) {
	if _, exists := dh.CandlesByTimeframe[timeframe]; !exists {
		dh.CandlesByTimeframe[timeframe] = make(map[time.Time]Candle)
	}
	for _, candle := range candles {
		dh.CandlesByTimeframe[timeframe][candle.Timestamp] = candle
	}
}

// GetCandle retrieves a candle for a specific timeframe and timestamp.
func (dh *DataHandler) GetCandle(timeframe Timeframe, timestamp time.Time) (Candle, error) {
	if candles, exists := dh.CandlesByTimeframe[timeframe]; exists {
		if candle, exists := candles[timestamp]; exists {
			return candle, nil
		}
	}
	return Candle{}, fmt.Errorf("no candle found for timeframe %s and timestamp %v", timeframe, timestamp)
}

// GetCandlesByTimeframe retrieves all candles for a specific timeframe.
func (dh *DataHandler) GetCandlesByTimeframe(timeframe Timeframe) ([]Candle, error) {
	if candles, exists := dh.CandlesByTimeframe[timeframe]; exists {
		var result []Candle
		for _, candle := range candles {
			result = append(result, candle)
		}
		return result, nil
	}
	return nil, fmt.Errorf("no candles found for timeframe %s", timeframe)
}

// Size returns the number of timeframes managed by the DataHandler.
func (dh *DataHandler) Size() int {
	return len(dh.CandlesByTimeframe)
}

// PreviousCandle retrieves the candle immediately before the given timestamp in the specified timeframe.
func (dh *DataHandler) PreviousCandle(timeframe Timeframe, timestamp time.Time) (Candle, error) {
	// Check if the timeframe exists in the data
	candlesMap, exists := dh.CandlesByTimeframe[timeframe]
	if !exists {
		return Candle{}, fmt.Errorf("no data available for timeframe %v", timeframe)
	}

	// Extract timestamps and sort them
	var timestamps []time.Time
	for ts := range candlesMap {
		timestamps = append(timestamps, ts)
	}
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	// Find the previous candle
	for i := 1; i < len(timestamps); i++ {
		if timestamps[i] == timestamp {
			return candlesMap[timestamps[i-1]], nil
		}
	}

	return Candle{}, fmt.Errorf("no previous candle found for timestamp %v in timeframe %v", timestamp, timeframe)
}
