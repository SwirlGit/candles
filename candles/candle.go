package candles

import (
	"fmt"
	"time"
)

type candle struct {
	ticker    string
	unixTime  time.Time
	maxPrice  float64
	minPrice  float64
	lastPrice float64
}

func newCandle(ticker string, t time.Time, price float64) *candle {
	return &candle{
		ticker:    ticker,
		unixTime:  t,
		maxPrice:  price,
		minPrice:  price,
		lastPrice: price,
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
	return fmt.Sprintf("[ticker: %v, unixTime: %v, maxPrice: %v, "+
		"minPrice: %v, lastPrice: %v]", c.ticker, c.unixTime,
		c.maxPrice, c.minPrice, c.lastPrice)
}

func (c *candle) ToCsvString() string {
	return fmt.Sprintf("%v;%v;%v;%v;%v", c.ticker, c.unixTime,
		c.maxPrice, c.minPrice, c.lastPrice)
}
