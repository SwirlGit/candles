package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-candles/candles"
)

func fanOut(in <-chan string) (<-chan string, <-chan string, <-chan string) {
	out1 := make(chan string)
	out2 := make(chan string)
	out3 := make(chan string)
	go func() {
		defer close(out1)
		defer close(out2)
		defer close(out3)
		for {
			line, ok := <-in
			if !ok {
				break
			}
			out1 <- line
			out2 <- line
			out3 <- line
		}
	}()
	return out1, out2, out3
}

func readFile(in <-chan struct{}, fileName string) (<-chan string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Can not open file %s: %s", fileName, err)
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer close(out)
		defer file.Close()
		<-in
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error while reading file %s: %s", fileName, err)
		}
		out <- "EOF"
	}()
	return out, nil
}

func parseLine(in <-chan string, duration time.Duration) <-chan string {
	handler := candles.NewHandler(duration)
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

func saveLine(in <-chan string, fileName string) (<-chan struct{}, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Can not open file %s: %s", fileName, err)
		return nil, err
	}

	out := make(chan struct{})
	go func() {
		defer close(out)
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
					line, fileName, err)
				continue
			}
			if n != len(line)+1 {
				log.Printf("Wrote %v symbols instead of %v while writing line %s",
					n, len(line)+1, line)
			}
		}
	}()
	return out, nil
}

func startPipeline(fileName string) (<-chan struct{}, error) {
	startChan := make(chan struct{})
	fileLinesOut, err := readFile(startChan, fileName)
	if err != nil {
		return nil, err
	}
	fileLinesOutShort, fileLinesOutMedium, fileLinesOutLong := fanOut(fileLinesOut)

	candlesOutShort := parseLine(fileLinesOutShort, 5*time.Minute)
	candlesOutMedium := parseLine(fileLinesOutMedium, 30*time.Minute)
	candlesOutLong := parseLine(fileLinesOutLong, 240*time.Minute)

	doneShort, err := saveLine(candlesOutShort, "candles_5min.csv")
	if err != nil {
		return nil, err
	}
	doneMedium, err := saveLine(candlesOutMedium, "candles_30min.csv")
	if err != nil {
		return nil, err
	}
	doneLong, err := saveLine(candlesOutLong, "candles_240min.csv")
	if err != nil {
		return nil, err
	}

	startChan <- struct{}{}
	close(startChan)

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-doneShort
		<-doneMedium
		<-doneLong
	}()

	return done, nil
}

func main() {
	fileNamePtr := flag.String("file", "", "path to file with trades")
	flag.Parse()
	if len(*fileNamePtr) == 0 {
		flag.Usage()
		return
	}

	done, err := startPipeline(*fileNamePtr)
	if err != nil {
		log.Printf("Can not start pipeline: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-done:
		return
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}
