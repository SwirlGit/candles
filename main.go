package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

func readFile(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		fileName, ok := <-in
		if ok {
			file, err := os.Open(fileName)
			if err != nil {
				// process open file error
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				out <- scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				// process scanner error
			}
		} else {
			// process in channel closing
		}
	}()
	return out
}

func parseLine(in <-chan string) <-chan string {
	var candle *Candle
	out := make(chan string)
	go func() {
		defer close(out)
		for line := range in {
			err := candle.Update(line)
			if err != nil {
				// process parsing error
			}
			out <- candle.String()
		}

	}()
	return out
}

func saveLine(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		file, err := os.Open("fileName")
		if err != nil {
			// process open file error
			return
		}
		defer file.Close()
		for line := range in {
			n, err := file.WriteString(line + "\n")
			if err != nil {
				// process err
			}
			if n != len(line) {
				// process wrong n number
			}
		}
	}()
	return out
}

func main() {
	// setup pipeline
	in := make(chan string)
	fileLinesOut := readFile(in)
	parsedLinesOut := parseLine(fileLinesOut)
	finalOut := saveLine(parsedLinesOut)

	<-finalOut
}
