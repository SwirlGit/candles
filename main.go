package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/candles/candles"
)

func readFile(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			fileName, ok := <-in
			if !ok {
				break
			}

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
			out <- "EOF"
		}
	}()
	return out
}

func parseLine(in <-chan string) <-chan string {
	handler := candles.NewHandler(5 * time.Minute)
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			line, ok := <-in
			if !ok {
				break
			}
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

func saveLine(in <-chan string) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		defer close(out)
		file, err := os.Create("candles_5min.csv")
		if err != nil {
			// process open file error
			return
		}
		defer file.Close()
		for {
			line, ok := <-in
			if !ok {
				break
			}
			if line == "EOF" {
				out <- struct{}{}
				continue
			}
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

func handlePipeline(fileName string) {
	filesIn := make(chan string)
	fileLinesOut := readFile(filesIn)
	candlesOut := parseLine(fileLinesOut)
	finalOut := saveLine(candlesOut)

	filesIn <- fileName
	<-finalOut
}

func main() {
	fileNamePtr := flag.String("file", "", "path to file with trades")
	flag.Parse()
	if len(*fileNamePtr) == 0 {
		flag.Usage()
		return
	}
	handlePipeline(*fileNamePtr)
}
