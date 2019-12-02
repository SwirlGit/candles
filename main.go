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
			log.Printf("Can not open file %s: %s", fileName, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error while reading file %s: %s", fileName, err)
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
				log.Printf("Error while parsing line %s: %s", line, err)
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
			log.Printf("Can not open file %s: %s", "candles_5min.csv", err)
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
				log.Printf("Error while writing line %s to file %s: %s",
					line, "candles_5min.csv", err)
				continue
			}
			if n != len(line) {
				log.Printf("Wrote %v symbols instead of %v while writing line %s",
					n, len(line), line)
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
