package kvcache

import (
	"backtesting/data"
	"time"
)

// KeyLevelType defines the type of key stored in the KeyStore.
type KeyLevelType string

const (
	LevelType       KeyLevelType = "level"
	RangeType       KeyLevelType = "range"
	CandlestickType KeyLevelType = "candlestick"
	FVGType         KeyLevelType = "fvg"
	MSBType         KeyLevelType = "msb"
	OrderBlockType  KeyLevelType = "orderBlock"
	SwingHighType   KeyLevelType = "swingHigh"
	SwingLowType    KeyLevelType = "swingLow"
)

// FVG represents a Fair Value Gap range.
type FVG struct {
	Low  float64
	High float64
}

// PriceRange represents a range with upper and lower bounds.
type PriceRange struct {
	Low  float64
	High float64
}
type SwingLow struct {
	Type string  //  "bearish"
	Low  float64 // The defining candlestick
}

type SwingHigh struct {
	Type string  //  "bearish"
	High float64 // The defining candlestick
}

// MSB represents a market structure break.
type MSB struct {
	Type       string  // "bullish" or "bearish"
	BreakPoint float64 // The price where the structure broke
}

// OrderBlock represents a bullish or bearish order block.
type OrderBlock struct {
	Type        string      // "bullish" or "bearish"
	Candlestick data.Candle // The defining candlestick
}

// KeyLevel represents a critical price level, range, candlestick, or concept like FVG/MSB/Order Block.
type KeyLevel struct {
	Name        string       // Unique identifier for the key
	Level       float64      // Single price level
	Range       *PriceRange  // Price range (for FVG or similar concepts)
	Candlestick *data.Candle // Important candlestick (optional)
	FVG         *FVG         // Fair Value Gap (optional)
	MSB         *MSB         // Market Structure Break (optional)
	OrderBlock  *OrderBlock  // Bullish/Bearish Order Block (optional)
	SwingLow    *SwingLow
	SwingHigh   *SwingHigh
	Type        KeyLevelType // Type of key
	Active      bool         // Whether the key is active
	Timestamp   time.Time    // Time when the key was identified
}

// KeyStore is a cache for tracking key price levels, ranges, and trading concepts.
type KeyStore struct {
	levels map[string]KeyLevel
}

// NewKeyStore initializes a new KeyStore.
func NewKeyStore() *KeyStore {
	return &KeyStore{
		levels: make(map[string]KeyLevel),
	}
}

// AddFVG adds a Fair Value Gap to the store.
func (ks *KeyStore) AddFVG(name string, low, high float64, candlestick data.Candle) {
	ks.levels[name] = KeyLevel{
		Name:      name,
		FVG:       &FVG{Low: low, High: high},
		Type:      FVGType,
		Active:    true,
		Timestamp: candlestick.Timestamp,
	}
}

// AddMSB adds a Market Structure Break to the store.
func (ks *KeyStore) AddMSB(name, msbType string, breakPoint float64, candlestick data.Candle) {
	ks.levels[name] = KeyLevel{
		Name:      name,
		MSB:       &MSB{Type: msbType, BreakPoint: breakPoint},
		Type:      MSBType,
		Active:    true,
		Timestamp: candlestick.Timestamp,
	}
}

// AddOrderBlock adds a Bullish/Bearish Order Block to the store.
func (ks *KeyStore) AddOrderBlock(name, blockType string, candlestick data.Candle) {
	ks.levels[name] = KeyLevel{
		Name:       name,
		OrderBlock: &OrderBlock{Type: blockType, Candlestick: candlestick},
		Type:       OrderBlockType,
		Active:     true,
		Timestamp:  candlestick.Timestamp,
	}
}

func (ks *KeyStore) AddSwingLow(name, blockType string, candlestick data.Candle) {
	ks.levels[name] = KeyLevel{
		Name:      name,
		SwingLow:  &SwingLow{Type: blockType, Low: candlestick.Low},
		Type:      SwingLowType,
		Active:    true,
		Timestamp: candlestick.Timestamp,
	}
}

func (ks *KeyStore) AddSwingHigh(name, blockType string, candlestick data.Candle) {
	ks.levels[name] = KeyLevel{
		Name:      name,
		SwingHigh: &SwingHigh{Type: blockType, High: candlestick.High},
		Type:      SwingHighType,
		Active:    true,
		Timestamp: candlestick.Timestamp,
	}
}

// CheckPriceBreach checks if a price breaches any active levels, FVGs, or MSBs.
func (ks *KeyStore) CheckPriceBreach(price float64) []KeyLevel {
	var breached []KeyLevel
	for name, key := range ks.levels {
		if !key.Active {
			continue
		}
		switch key.Type {
		case LevelType:
			if price <= key.Level {
				breached = append(breached, key)
				key.Active = false
				ks.levels[name] = key
			}
		case RangeType, FVGType:
			if key.Range != nil && price >= key.Range.Low && price <= key.Range.High {
				breached = append(breached, key)
			}
		case MSBType:
			if key.MSB != nil && ((key.MSB.Type == "bullish" && price > key.MSB.BreakPoint) ||
				(key.MSB.Type == "bearish" && price < key.MSB.BreakPoint)) {
				breached = append(breached, key)
			}
		}
	}
	return breached
}
