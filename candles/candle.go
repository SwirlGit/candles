package candles

import (
	"fmt"
	"time"
)

type candle struct {
	ticker     string
	unixTime   time.Time
	maxPrice   float64
	minPrice   float64
	firstPrice float64
	lastPrice  float64
}

func newCandle(ticker string, t time.Time, price float64) *candle {
	return &candle{
		ticker:     ticker,
		unixTime:   t,
		maxPrice:   price,
		minPrice:   price,
		firstPrice: price,
		lastPrice:  price,
	}
}

func (c *candle) updatePrice(price float64) {
	if price > c.maxPrice {
		c.maxPrice = price
	}
	if price < c.minPrice {
		c.minPrice = price
	}
	c.lastPrice = price
}

func (c *candle) String() string {
	return fmt.Sprintf("[ticker: %v, unixTime: %v, firstPrice: %v, "+
		"maxPrice: %v, minPrice: %v, lastPrice: %v]", c.ticker,
		c.unixTime.Format(time.RFC3339), c.maxPrice, c.minPrice,
		c.firstPrice, c.lastPrice)
}

func (c *candle) ToCsvString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v", c.ticker,
		c.unixTime.Format(time.RFC3339), c.firstPrice,
		c.maxPrice, c.minPrice, c.lastPrice)
}
