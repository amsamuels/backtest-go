package strategies

import (
	data "backtesting/data"
	"backtesting/kvcache"
	log "backtesting/logger"
)

// ICTStrategy implements the ICT Smart Money Concept strategy.
type ICTStrategy struct {
	keyStore *kvcache.KeyStore
}
type Trade = data.Trade
type Candle = data.Candle

// NewICTStrategy initializes a new ICTStrategy.
func NewICTStrategy(store *kvcache.KeyStore) *ICTStrategy {
	return &ICTStrategy{
		keyStore: store,
	}
}
func (s *ICTStrategy) Run(data *data.DataHandler, timeframe data.Timeframe) ([]Trade, error) {
	log.Info("Starting ICT Smart Money strategy for timeframe: %v", timeframe)

	var trades []Trade

	// Fetch candles for the specified timeframe
	candles, err := data.GetCandlesByTimeframe(timeframe)
	if err != nil {
		log.Error("Error fetching candles for timeframe %v: %v", timeframe, err)
		return nil, err
	}

	// Iterate through candles
	for i, candle := range candles {
		var prevCandle *Candle
		var nextCandle *Candle

		if i > 0 {
			prevCandle = &candles[i-1]
		}
		if i < len(candles)-1 {
			nextCandle = &candles[i+1]
		}

		// Check for breaches in the key store
		breachedLevels := s.keyStore.CheckPriceBreach(candle.Low)
		for _, level := range breachedLevels {
			log.Info("Breached key level: %s at price %.2f", level.Type, level.Level)
			trades = append(trades, Trade{
				Timestamp: candle.Timestamp,
				Type:      "Sell",
				Price:     candle.Low,
			})
		}

		//	Add or update swing lows/highs dynamically
		if prevCandle != nil && nextCandle != nil && candle.IsSwingLow(*prevCandle, *nextCandle) {
			// TODO: FIND A WAY TO MAKE EACH STORED ENTRY UNIQUE
			s.keyStore.AddSwingLow("", "swingLow", candle)
		}
		if prevCandle != nil && nextCandle != nil && candle.IsSwingHigh(*prevCandle, *nextCandle) {
			s.keyStore.AddSwingLow("", "swingHigh", candle)
		}
	}

	log.Info("ICT Smart Money strategy completed for timeframe: %v", timeframe)
	return trades, nil
}
