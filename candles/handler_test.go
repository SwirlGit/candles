package candles

import (
	"testing"
	"time"
)

func TestDayStart(t *testing.T) {
	type DayStartTestCase struct {
		givenTime time.Time
		result    time.Time
	}
	tables := []DayStartTestCase{
		DayStartTestCase{
			givenTime: time.Date(2019, 12, 1, 12, 37, 45, 888, time.UTC),
			result:    time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		DayStartTestCase{
			givenTime: time.Date(2019, 12, 1, 23, 59, 59, 888, time.UTC),
			result:    time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, table := range tables {
		resultDayStart := dayStart(table.givenTime)
		if table.result != resultDayStart {
			t.Fatalf("TestDayStart: wrong day start, expected: %s, got: %s",
				table.result, resultDayStart)
		}
	}
}

func TestValidateInputValues(t *testing.T) {
	type ValidateInputValuesTestCase struct {
		values inputValues
		result bool
	}
	tables := []ValidateInputValuesTestCase{
		ValidateInputValuesTestCase{
			values: inputValues{
				unixTime: time.Date(2019, 12, 1, 7, 0, 0, 0, time.UTC),
			},
			result: true,
		},
		ValidateInputValuesTestCase{
			values: inputValues{
				unixTime: time.Date(2019, 12, 1, 23, 59, 59, 999999999, time.UTC),
			},
			result: true,
		},
		ValidateInputValuesTestCase{
			values: inputValues{
				unixTime: time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			result: false,
		},
		ValidateInputValuesTestCase{
			values: inputValues{
				unixTime: time.Date(2019, 12, 1, 6, 59, 59, 999999999, time.UTC),
			},
			result: false,
		},
	}
	for _, table := range tables {
		validationResult := validateInputValues(table.values)
		if table.result != validationResult {
			t.Fatalf("TestValidateInputValues: wrong validation result for time %s, "+
				"expected: %v, got: %v", table.values.unixTime, table.result, validationResult)
		}
	}
}

func TestBaseTime(t *testing.T) {
	type BaseTimeTestCase struct {
		givenTime time.Time
		duration  time.Duration
		result    time.Time
	}
	tables := []BaseTimeTestCase{
		BaseTimeTestCase{
			givenTime: time.Date(2019, 12, 1, 12, 37, 45, 888, time.UTC),
			duration:  5 * time.Minute,
			result:    time.Date(2019, 12, 1, 12, 35, 0, 0, time.UTC),
		},
		BaseTimeTestCase{
			givenTime: time.Date(2019, 12, 1, 1, 12, 59, 888, time.UTC),
			duration:  14 * time.Minute,
			result:    time.Date(2019, 12, 1, 1, 10, 0, 0, time.UTC),
		},
	}
	for _, table := range tables {
		resultBaseTime := baseTime(table.givenTime, table.duration)
		if table.result != resultBaseTime {
			t.Fatalf("TestBaseTime: wrong base time, expected: %s, got: %s",
				table.result, resultBaseTime)
		}
	}
}

func TestNewHandler(t *testing.T) {
	type NewHandlerTestCase struct {
		duration time.Duration
		result   Handler
	}
	tables := []NewHandlerTestCase{
		NewHandlerTestCase{
			duration: 5 * time.Minute,
			result: Handler{
				duration: 5 * time.Minute,
			},
		},
	}
	for _, table := range tables {
		resultHandler := NewHandler(table.duration)
		if table.result.duration != resultHandler.duration {
			t.Fatalf("TestNewHandler: wrong time duration, expected: %s, got: %s",
				table.result.duration, resultHandler.duration)
		}
	}
}

func TestProcessLine(t *testing.T) {
	type ProcessLineTestCase struct {
		handler  *Handler
		newLines []string
		result   [][]string
	}
	tables := []ProcessLineTestCase{
		ProcessLineTestCase{
			handler: &Handler{
				candles:  make(map[string]*candle),
				duration: 1 * time.Minute,
			},
			newLines: []string{
				"AMZN,1645,3,2019-01-31 07:00:01.970000",
				"SBR,250.67,3,2019-01-31 07:00:01.980000",
				"AMZN,1675.6,3,2019-01-31 07:00:30.970000",
				"SBR,256.67,3,2019-01-31 07:00:30.980000",
				"AMZN,1675,3,2019-01-31 07:01:00.000000",
				"SBR,258.67,3,2019-01-31 07:01:30.000000",
			},
			result: [][]string{
				{},
				{},
				{},
				{},
				{"AMZN,2019-01-31T07:00:00Z,1645,1675.6,1645,1675",
					"SBR,2019-01-31T07:00:00Z,250.67,258.67,250.67,258.67"},
				{},
			},
		},
	}
	for _, table := range tables {
		for i := 0; i < len(table.newLines); i++ {
			resultStrings, err := table.handler.ProcessLine(table.newLines[i])
			if err != nil {
				t.Fatalf("TestProcessLine: failed to process new line: %s", err)
			}
			if len(table.result[i]) != len(resultStrings) {
				t.Fatalf("TestProcessLine: wrong number of outcome strings, "+
					"expected: %v, got: %v", len(table.result[i]), len(resultStrings))
			}
			for j := 0; j < len(table.result[j]); j++ {
				if table.result[i][j] != resultStrings[j] {
					t.Fatalf("TestProcessLine: wrong out string step: %v, line: %v, "+
						"expected: %v, got: %v", i+1, j+1, table.result[i][j], resultStrings[j])
				}
			}
		}
	}
}

func TestClose(t *testing.T) {
	type CloseTestCase struct {
		handler *Handler
		result  []string
	}
	tables := []CloseTestCase{
		CloseTestCase{
			handler: &Handler{
				candles:  make(map[string]*candle),
				duration: 1 * time.Minute,
			},
			result: []string{},
		},
		CloseTestCase{
			handler: &Handler{
				candles: map[string]*candle{
					"Ticker1": &candle{
						ticker:     "Ticker1",
						unixTime:   time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
						maxPrice:   5.0,
						minPrice:   1.0,
						firstPrice: 2.0,
						lastPrice:  4.0,
					},
					"Ticker2": &candle{
						ticker:     "Ticker2",
						unixTime:   time.Date(2, 2, 2, 2, 2, 2, 2, time.UTC),
						maxPrice:   15.75,
						minPrice:   10.6,
						firstPrice: 11.0,
						lastPrice:  12.0,
					},
				},
				duration: 1 * time.Minute,
			},
			result: []string{"Ticker1,0001-01-01T01:01:01Z,2,5,1,4",
				"Ticker2,0002-02-02T02:02:02Z,11,15.75,10.6,12",
			},
		},
	}
	for _, table := range tables {
		resultStrings := table.handler.Close()
		if len(table.result) != len(resultStrings) {
			t.Fatalf("TestClose: wrong number of outcome strings, "+
				"expected: %v, got: %v", len(table.result), len(resultStrings))
		}
		for i := 0; i < len(table.result); i++ {
			if table.result[i] != resultStrings[i] {
				t.Fatalf("TestClose: wrong out string line: %v, "+
					"expected: %v, got: %v", i+1, table.result[i], resultStrings[i])
			}
		}
	}
}
