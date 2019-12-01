package main

import (
	"bufio"
	"os"

	"github.com/candles/candles"
)

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
	var candle *candles.Candle
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
