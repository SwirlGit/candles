package main

import (
	"bufio"
	"os"
	"time"

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
	handler := candles.NewHandler(5 * time.Minute)
	out := make(chan string)
	go func() {
		defer close(out)
		for line := range in {
			outStrings, err := handler.ProcessLine(line)
			if err != nil {
				// process parsing error
				continue
			}
			for _, outString := range outStrings {
				out <- outString
			}
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
