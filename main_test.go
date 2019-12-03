package main

import (
	"bufio"
	"os"
	"testing"
	"time"
)

// проверка текущего прочитанного блока
func checkCurrentCandles(t *testing.T, testCandles, resultCandles map[string]string, fileName string) {
	for testCandleTicker, testCandle := range testCandles {
		found := false
		for resultCandleTicker, resultCandle := range resultCandles {
			if testCandle == resultCandle {
				delete(testCandles, testCandleTicker)
				delete(resultCandles, resultCandleTicker)
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("TestHandlePipeline: have not found in result file %s test "+
				"line: %s", fileName, testCandle)
		}
	}
	for _, resultCandle := range resultCandles {
		t.Fatalf("TestHandlePipeline: found extra candle: %s", resultCandle)
	}
}

// так как порядок выходных строк не фиксирован, кроме
// сортировки по времени, то будем проверять данные блоками, где
// в каждом блоке все данные для одного времени открытия
func checkFile(t *testing.T, testDataFileName, resultDataFileName string) {
	testDataFile, err := os.Open(testDataFileName)
	if err != nil {
		t.Fatalf("TestHandlePipeline: can not open test data file: %s", err)
	}

	resultDataFile, err := os.Open(resultDataFileName)
	if err != nil {
		t.Fatalf("TestHandlePipeline: can not open result data file: %s", err)
	}

	testDataScanner := bufio.NewScanner(testDataFile)
	testCandles := make(map[string]string)
	resultDataScanner := bufio.NewScanner(resultDataFile)
	resultCandles := make(map[string]string)

	var currentTime time.Time
	isFirstLine := true
	for testDataScanner.Scan() {
		newTime, err := time.Parse(time.RFC3339, testDataScanner.Text()[5:25])
		if err != nil {
			t.Fatalf("can not parse time: %s", err)
		}
		if !isFirstLine && currentTime != newTime {
			checkCurrentCandles(t, testCandles, resultCandles, testDataFileName)
		}
		isFirstLine = false
		currentTime = newTime

		testCandles[testDataScanner.Text()[0:4]] = testDataScanner.Text()
		if !resultDataScanner.Scan() {
			t.Fatal("TestHandlePipeline: result data file has less lines than expected")
		}
		resultCandles[resultDataScanner.Text()[0:4]] = resultDataScanner.Text()
	}
	if !isFirstLine {
		checkCurrentCandles(t, testCandles, resultCandles, testDataFileName)
	}
	if resultDataScanner.Scan() {
		t.Fatal("TestHandlePipeline: result data file has more lines than expected")
	}
}

func TestStartPipeline(t *testing.T) {
	done, err := startPipeline("test/trades.csv")
	if err != nil {
		t.Fatalf("TestStartPipeline: can not start pipeline: %s", err)
	}
	<-done

	checkFile(t, "test/candles_5min.csv", "candles_5min.csv")
	checkFile(t, "test/candles_30min.csv", "candles_30min.csv")
	checkFile(t, "test/candles_240min.csv", "candles_240min.csv")
}
