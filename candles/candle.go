package candles

import (
	"fmt"
	"time"
)

type candle struct {
	ticker    string
	unixTime  time.Time
	maxPrize  float64
	minPrize  float64
	lastPrize float64
}

func createCandle(ticker string, t time.Time, prize float64) candle {
	return candle{
		ticker:    ticker,
		unixTime:  t,
		maxPrize:  prize,
		minPrize:  prize,
		lastPrize: prize,
	}
}

func (c candle) updatePrize(prize float64) {
	if prize > c.maxPrize {
		c.maxPrize = prize
	}
	if prize < c.minPrize {
		c.minPrize = prize
	}
	c.lastPrize = prize
}

func (c candle) String() string {
	return fmt.Sprintf("%v;%v;%v;%v;%v", c.ticker, c.unixTime,
		c.maxPrize, c.minPrize, c.lastPrize)
}
