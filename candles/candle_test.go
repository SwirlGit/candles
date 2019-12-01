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
				ticker:     "Ticker",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   150.7,
				minPrice:   150.7,
				firstPrice: 150.7,
				lastPrice:  150.7,
			},
		},
	}
	for _, table := range tables {
		newCandleResult := newCandle(table.ticker, table.unixTime,
			table.price)
		if table.result != *newCandleResult {
			t.Fatalf("TestNewCandle: wrong candle, expected: %s, got: %s",
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
				ticker:     "Ticker",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   90.0,
				minPrice:   90.0,
				firstPrice: 90.0,
				lastPrice:  90.0,
			},
			newPrices: []float64{1.0, 100.0, 150.0},
			result: candle{
				ticker:     "Ticker",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   150.0,
				minPrice:   1.0,
				firstPrice: 90.0,
				lastPrice:  150.0,
			},
		},
		UpdatePriceTestCase{
			initial: candle{
				ticker:     "Ticker",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   90.0,
				minPrice:   90.0,
				firstPrice: 90.0,
				lastPrice:  90.0,
			},
			newPrices: []float64{100, 150.0, 125.0},
			result: candle{
				ticker:     "Ticker",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   150.0,
				minPrice:   90.0,
				firstPrice: 90.0,
				lastPrice:  125.0,
			},
		},
	}
	for _, table := range tables {
		candleToUpdate := table.initial
		for i := 0; i < len(table.newPrices); i++ {
			candleToUpdate.updatePrice(table.newPrices[i])
		}
		if table.result != candleToUpdate {
			t.Fatalf("TestUpdatePrice: wrong candle, expected: %s, got: %s",
				&table.result, &candleToUpdate)
		}
	}
}

func TestToCsvString(t *testing.T) {
	type ToCsvStringTestCase struct {
		candleToConver candle
		result         string
	}
	tables := []ToCsvStringTestCase{
		ToCsvStringTestCase{
			candleToConver: candle{
				ticker:     "Ticker1",
				unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
				maxPrice:   150.11,
				minPrice:   90.0,
				firstPrice: 97.7,
				lastPrice:  125.0,
			},
			result: "Ticker1,0001-01-01T01:01:01Z,97.7,150.11,90,125",
		},
	}
	for _, table := range tables {
		resultCsvString := table.candleToConver.ToCsvString()
		if table.result != resultCsvString {
			t.Fatalf("TestToCsvString: wrong csv string, expected: %s, got: %s",
				table.result, resultCsvString)
		}
	}
}
