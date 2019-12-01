package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/candles/candles"
)

func readFile(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		fileName := <-in
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
				return
			}
			if line == "EOF" {
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

func startPipeline(fileName string) <-chan struct{} {
	filesIn := make(chan string)
	fileLinesOut := readFile(filesIn)
	candlesOut := parseLine(fileLinesOut)
	done := saveLine(candlesOut)

	filesIn <- fileName
	close(filesIn)

	return done
}

func main() {
	fileNamePtr := flag.String("file", "", "path to file with trades")
	flag.Parse()
	if len(*fileNamePtr) == 0 {
		flag.Usage()
		return
	}

	done := startPipeline(*fileNamePtr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-done:
		return
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}
