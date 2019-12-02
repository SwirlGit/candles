package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-candles/candles"
)

func fanIn(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func fanOut(in <-chan string) (<-chan string, <-chan string, <-chan string) {
	out1 := make(chan string)
	out2 := make(chan string)
	out3 := make(chan string)
	go func() {
		defer close(out1)
		defer close(out2)
		defer close(out3)
		for n := range in {
			out1 <- n
			out2 <- n
			out3 <- n
		}
	}()
	return out1, out2, out3
}

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

func saveLine(in <-chan string, fileName string) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		defer close(out)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("Can not open file %s: %s", fileName, err)
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
					line, fileName, err)
				continue
			}
			if n != len(line)+1 {
				log.Printf("Wrote %v symbols instead of %v while writing line %s",
					n, len(line)+1, line)
			}
		}
	}()
	return out
}

func startPipeline(fileName string) <-chan struct{} {
	filesIn := make(chan string)
	fileLinesOut := readFile(filesIn)
	fileLinesOutShort, fileLinesOutMedium, fileLinesOutLong := fanOut(fileLinesOut)

	candlesOutShort := parseLine(fileLinesOutShort, 5*time.Minute)
	candlesOutMedium := parseLine(fileLinesOutMedium, 30*time.Minute)
	candlesOutLong := parseLine(fileLinesOutLong, 240*time.Minute)

	doneShort := saveLine(candlesOutShort, "candles_5min.csv")
	doneMedium := saveLine(candlesOutMedium, "candles_30min.csv")
	doneLong := saveLine(candlesOutLong, "candles_240min.csv")

	filesIn <- fileName
	close(filesIn)

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-doneShort
		<-doneMedium
		<-doneLong
	}()

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
