package main

import (
	"backtesting/data"
	log "backtesting/logger"
	"fmt"
	"os"
)

func main() {

	// Initialize the logger
	err := log.Init("program.log")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Close()

	// Map of file paths to timeframes
	files := map[string]data.Timeframe{
		"Bitstamp_ETHUSDT_d.csv": data.Daily,
	}

	// Load candles for multiple timeframes
	dataHandler, err := data.LoadCandlesByTimeframes(files)
	if err != nil {
		log.Error("Error loading candles: %v", err)
		log.Close() // Flush logs before exiting
		os.Exit(1)
	}

	// List all daily candles
	candles, err := dataHandler.GetCandlesByTimeframe(data.Daily)
	if err != nil {
		log.Error("Error loading candles: %v", err)
	}

	log.Info("Trade executed: %+v", candles)
	// // Load data
	// dataHandler := &data.DataHandler{
	// 	CandlesByTimeframe: candles,
	// }

	// // Initialize strategy
	// strategy := strategies.NewICTStrategy()

	// // Run strategy for 1-minute timeframe
	// trades, err := strategy.Run(dataHandler, data.OneMinute)
	// if err != nil {
	// 	log.Error("Error running strategy: %v", err)
	// 	return
	// }

	// Print trades
	// for _, trade := range trades {
	// 	log.Info("Trade executed: %+v", trade)
	// }

}
