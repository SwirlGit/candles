package candles

import (
	"errors"
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
var errWrongUnixTime = errors.New("wrong unix time")
var validStartDuration = 7 * time.Hour

type inputValues struct {
	ticker   string
	unixTime time.Time
	price    float64
}

func parseInputLine(line string) (inputValues, error) {
	var values inputValues
	inputStrings := strings.Split(line, sep)
	if len(inputStrings) != numOfValues {
		return values, errWrongNumberOfParameters
	}
	ticker := inputStrings[0]
	price, err := strconv.ParseFloat(inputStrings[1], 64)
	if err != nil {
		return values, err
	}
	t, err := time.Parse(timeFormat, inputStrings[2])
	if err != nil {
		return values, err
	}
	t = t.UTC()

	values = inputValues{
		ticker:   ticker,
		unixTime: t,
		price:    price,
	}
	return values, nil
}

func dayStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func validateInputValues(values inputValues) bool {
	start := dayStart(values.unixTime)
	return values.unixTime.Sub(start) >= validStartDuration
}

func baseTime(t time.Time, d time.Duration) time.Time {
	start := dayStart(t)
	duration := t.Sub(start)
	index := duration / d
	return start.Add(index * d)
}

// Handler управляющий набором свечей
type Handler struct {
	candles  map[string]candle
	duration time.Duration
	lastTime time.Time
}

// NewHandler функция конструктор для обработчика
// timeInterval в минутах
func NewHandler(timeInterval int) *Handler {
	return &Handler{
		duration: time.Duration(timeInterval) * time.Minute,
	}
}

// ProcessLine обработать строку
func (handler *Handler) ProcessLine(line string) ([]string, error) {
	values, err := parseInputLine(line)
	if err != nil {
		return []string{}, err
	}

	if !validateInputValues(values) {
		return []string{}, errWrongUnixTime
	}

	var candlesStrings []string
	for ticker, currentCandle := range handler.candles {
		if values.unixTime.Sub(currentCandle.unixTime) > handler.duration {
			candlesStrings = append(candlesStrings, currentCandle.String())
			delete(handler.candles, ticker)
		}
	}

	if _, exist := handler.candles[values.ticker]; !exist {
		handler.candles[values.ticker] = createCandle(values.ticker,
			baseTime(values.unixTime, handler.duration), values.price)
	} else {
		handler.candles[values.ticker].updatePrice(values.price)
	}

	return candlesStrings, nil
}

// Close закончить обработку
func (handler *Handler) Close() []string {
	var candlesStrings []string
	for ticker, currentCandle := range handler.candles {
		candlesStrings = append(candlesStrings, currentCandle.String())
		delete(handler.candles, ticker)
	}
	return candlesStrings
}
