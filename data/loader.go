package data

import (
	log "backtesting/logger"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

var requiredColumns = []string{"timestamp", "open", "high", "low", "close", "volume"}

// LoadCandlesByTimeframes loads candles for multiple timeframes from a set of CSV files.
func LoadCandlesByTimeframes(files map[string]Timeframe) (*DataHandler, error) {
	dh := NewDataHandler()
	for filePath, timeframe := range files {
		candles, err := loadCandlesFromFile(filePath)
		if err != nil {
			log.Error("Error loading candles: %v", err)
			return nil, fmt.Errorf("error loading candles from file %s: %w", filePath, err)
		}
		dh.AddCandles(timeframe, candles)
	}

	return dh, nil
}

// loadCandlesFromFile loads candles from a single CSV file.
func loadCandlesFromFile(csvFilePath string) ([]Candle, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Error("error opening file: %s", err)
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and validate header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	columnMap, err := mapColumns(header)
	if err != nil {
		return nil, err
	}

	// Parse candles
	var candles []Candle

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		// Parse a single candle
		candle, err := parseRow(row, columnMap)
		if err != nil {
			// Skip invalid rows
			continue
		}

		candles = append(candles, candle)
	}

	return candles, nil
}

// mapColumns maps the header to the required fields.
func mapColumns(header []string) (map[string]int, error) {
	columnMap := make(map[string]int)

	// Match required columns
	for i, col := range header {
		switch col {
		case "timestamp", "unix":
			columnMap["timestamp"] = i
		case "open":
			columnMap["open"] = i
		case "high":
			columnMap["high"] = i
		case "low":
			columnMap["low"] = i
		case "close":
			columnMap["close"] = i
		case "volume", "Volume ETH", "Volume USDT":
			columnMap["volume"] = i
		}
	}

	// Ensure all required fields are mapped
	for _, field := range requiredColumns {
		if _, exists := columnMap[field]; !exists {
			return nil, fmt.Errorf("missing required column: %s", field)
		}
	}

	return columnMap, nil
}

// parseRow parses a single row into a Candle using the column mapping.
func parseRow(row []string, columnMap map[string]int) (Candle, error) {
	var candle Candle
	var err error

	// Parse timestamp
	if idx, exists := columnMap["timestamp"]; exists {
		timestampStr := row[idx]
		candle.Timestamp, err = parseUnixTimestamp(timestampStr)
		if err != nil {
			return Candle{}, fmt.Errorf("invalid timestamp: %v", err)
		}
	}

	// Parse numerical fields
	if idx, exists := columnMap["open"]; exists {
		candle.Open, _ = strconv.ParseFloat(row[idx], 64)
	}
	if idx, exists := columnMap["high"]; exists {
		candle.High, _ = strconv.ParseFloat(row[idx], 64)
	}
	if idx, exists := columnMap["low"]; exists {
		candle.Low, _ = strconv.ParseFloat(row[idx], 64)
	}
	if idx, exists := columnMap["close"]; exists {
		candle.Close, _ = strconv.ParseFloat(row[idx], 64)
	}
	if idx, exists := columnMap["volume"]; exists {
		candle.Volume, _ = strconv.ParseFloat(row[idx], 64)
	}

	return candle, nil
}

// parseUnixTimestamp parses a Unix timestamp string into a time.Time object.
func parseUnixTimestamp(ts string) (time.Time, error) {
	unixTime, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing Unix timestamp: %v", err)
	}
	return time.Unix(unixTime, 0), nil
}
