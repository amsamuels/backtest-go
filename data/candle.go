package data

import "time"

// Candle represents OHLCV data for a specific time interval.
type Candle struct {
	Timestamp time.Time // Start time of the candle
	Open      float64   // Opening price
	High      float64   // Highest price during the interval
	Low       float64   // Lowest price during the interval
	Close     float64   // Closing price
	Volume    float64   // Volume traded during the interval
}

// NewCandle creates a new Candle from raw OHLCV data.
func NewCandle(timestamp time.Time, open, high, low, close, volume float64) Candle {
	return Candle{
		Timestamp: timestamp,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}

// Range returns the range of the candle (High - Low).
func (c *Candle) Range() float64 {
	return c.High - c.Low
}

// BodySize returns the size of the candle body (Close - Open).
func (c *Candle) BodySize() float64 {
	return c.Close - c.Open
}

// IsBullish checks if the candle is bullish (Close > Open).
func (c *Candle) IsBullish() bool {
	return c.Close > c.Open
}

// IsBearish checks if the candle is bearish (Close < Open).
func (c *Candle) IsBearish() bool {
	return c.Close < c.Open
}

// IsEqualHigh checks if two candles have similar highs within a tolerance.
func (c *Candle) IsEqualHigh(other Candle, tolerance float64) bool {
	return abs(c.High-other.High) <= tolerance
}

// IsEqualLow checks if two candles have similar lows within a tolerance.
func (c *Candle) IsEqualLow(other Candle, tolerance float64) bool {
	return abs(c.Low-other.Low) <= tolerance
}

// SwingHighLiquidityZone returns the high of a swing high as a potential liquidity zone.
func (c *Candle) SwingHighLiquidityZone() float64 {
	return c.High
}

// SwingLowLiquidityZone returns the low of a swing low as a potential liquidity zone.
func (c *Candle) SwingLowLiquidityZone() float64 {
	return c.Low
}

// Helper function: abs calculates the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// IsBullishOrderBlock identifies a bullish order block.
func (c *Candle) IsBullishOrderBlock() bool {
	return c.IsBearish() && c.BodySize() > (c.Range()*0.5) // Body > 50% of range
}

// IsBearishOrderBlock identifies a bearish order block.
func (c *Candle) IsBearishOrderBlock() bool {
	return c.IsBullish() && c.BodySize() > (c.Range()*0.5)
}

// OrderBlockLevel returns the level of the order block (open price).
func (c *Candle) OrderBlockLevel() float64 {
	return c.Open
}

// HasFairValueGap checks if there's a gap between two candles.
func (c *Candle) HasFairValueGap(next Candle) bool {
	return c.Low > next.High || next.Low > c.High
}

// FVGRange calculates the range of a Fair Value Gap (if present).
func (c *Candle) FVGRange(next Candle) (float64, float64, bool) {
	if c.HasFairValueGap(next) {
		if c.Low > next.High {
			return next.High, c.Low, true
		} else {
			return c.High, next.Low, true
		}
	}
	return 0, 0, false
}

// IsSwingHigh checks if the current candle is a swing high.
func (c *Candle) IsSwingHigh(previous, next Candle) bool {
	return c.High > previous.High && c.High > next.High
}

// IsSwingLow checks if the current candle is a swing low.
func (c *Candle) IsSwingLow(previous, next Candle) bool {
	return c.Low < previous.Low && c.Low < next.Low
}

// IsBreakOfStructure checks if a candle breaks a given high or low.
func (c *Candle) IsBreakOfStructure(level float64, direction string) bool {
	if direction == "up" {
		return c.High > level
	} else if direction == "down" {
		return c.Low < level
	}
	return false
}
