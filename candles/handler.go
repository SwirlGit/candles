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

// TODO: проверка на закрытие свечи

type inputValues struct {
	ticker   string
	unixTime time.Time
	prize    float64
}

func parseInputLine(line string) (inputValues, error) {
	var values inputValues
	inputStrings := strings.Split(line, sep)
	if len(inputStrings) != numOfValues {
		return values, errWrongNumberOfParameters
	}
	ticker := inputStrings[0]
	prize, err := strconv.ParseFloat(inputStrings[1], 64)
	if err != nil {
		return values, err
	}
	t, err := time.Parse(timeFormat, inputStrings[2])
	if err != nil {
		return values, err
	}
	values = inputValues{
		ticker:   ticker,
		unixTime: t,
		prize:    prize,
	}
	return values, nil
}

// Handler управляющий набором свечей
type Handler struct {
	candles      map[string]candle
	timeInterval int
	lastTime     time.Time
}

// NewHandler функция конструктор для обработчика
func NewHandler(timeInterval int) *Handler {
	return &Handler{
		timeInterval: timeInterval,
	}
}

// ProcessNewLine обработать строку
func (handler *Handler) ProcessLine(line string) ([]string, error) {
	_, err := parseInputLine(line)
	if err != nil {
		return []string{}, err
	}
	return []string{}, nil
}

// Close закончить обработку
func (handler *Handler) Close() []string {
	return []string{}
}
