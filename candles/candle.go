package candles

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	sep         = ";"
	numOfValues = 4
	timeFormat  = "2006-01-02 15:04:05.000006"
)

var errWrongNumberOfParameters = errors.New("wrong number of parameters")

// TODO: проверка на закрытие свечи

// Candle структура свечи
type Candle struct {
	ticker    string
	unixTime  time.Time
	maxPrize  float64
	minPrize  float64
	lastPrize float64
}

// NewCandle функция конструктор свечи
func NewCandle(ticker string, t time.Time, prize float64) *Candle {
	return &Candle{
		ticker:    ticker,
		unixTime:  t,
		maxPrize:  prize,
		minPrize:  prize,
		lastPrize: prize,
	}
}

func (candle *Candle) updatePrize(prize float64) {
	if prize > candle.maxPrize {
		candle.maxPrize = prize
	}
	if prize < candle.minPrize {
		candle.minPrize = prize
	}
	candle.lastPrize = prize
}

// Update обновить значения свечи новыми данными
func (candle *Candle) Update(csv string) error {
	incomeValues := strings.Split(csv, sep)
	if len(incomeValues) != numOfValues {
		return errWrongNumberOfParameters
	}
	ticker := incomeValues[0]
	prize, err := strconv.ParseFloat(incomeValues[1], 64)
	if err != nil {
		return err
	}
	t, err := time.Parse(timeFormat, incomeValues[2])
	if err != nil {
		return err
	}
	if candle == nil {
		candle = NewCandle(ticker, t, prize)
	} else {
		candle.updatePrize(prize)
	}
	return nil
}

func (candle *Candle) String() string {
	return fmt.Sprintf("%v;%v;%v;%v;%v", candle.ticker, candle.unixTime,
		candle.maxPrize, candle.minPrize, candle.lastPrize)
}
