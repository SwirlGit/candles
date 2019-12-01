package candles

import (
	"testing"
	"time"
)

func TestNewCandle(t *testing.T) {
	type NewCandleTestCase struct {
		ticker   string
		unixTime time.Time
		price    float64
		result   candle
	}
	tables := []NewCandleTestCase{
		NewCandleTestCase{
			ticker:   "Ticker",
			unixTime: time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
			price:    150.7,
			result: candle{
				ticker:    "Ticker",
				unixTime:  time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:  150.7,
				minPrice:  150.7,
				lastPrice: 150.7,
			},
		},
	}
	for _, table := range tables {
		newCandleResult := newCandle(table.ticker, table.unixTime,
			table.price)
		if table.result != *newCandleResult {
			t.Fatalf("newCandle: wrong candle, expected: %s, got: %s",
				&table.result, newCandleResult)
		}
	}
}

func TestUpdatePrice(t *testing.T) {
	type UpdatePriceTestCase struct {
		initial   candle
		newPrices []float64
		result    candle
	}
	tables := []UpdatePriceTestCase{
		UpdatePriceTestCase{
			initial: candle{
				ticker:    "Ticker",
				unixTime:  time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:  90.0,
				minPrice:  90.0,
				lastPrice: 90.0,
			},
			newPrices: []float64{1.0, 100.0, 150.0},
			result: candle{
				ticker:    "Ticker",
				unixTime:  time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:  150.0,
				minPrice:  1.0,
				lastPrice: 150.0,
			},
		},
		UpdatePriceTestCase{
			initial: candle{
				ticker:    "Ticker",
				unixTime:  time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:  90.0,
				minPrice:  90.0,
				lastPrice: 90.0,
			},
			newPrices: []float64{100, 150.0, 125.0},
			result: candle{
				ticker:    "Ticker",
				unixTime:  time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:  150.0,
				minPrice:  90.0,
				lastPrice: 125.0,
			},
		},
	}
	for _, table := range tables {
		candleToUpdate := table.initial
		for i := 0; i < len(table.newPrices); i++ {
			candleToUpdate.updatePrice(table.newPrices[i])
		}
		if table.result != candleToUpdate {
			t.Fatalf("updatePrice: wrong candle, expected: %s, got: %s",
				&table.result, &candleToUpdate)
		}
	}
}
