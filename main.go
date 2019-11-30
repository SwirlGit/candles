package main

import (
	"bufio"
	"os"
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
	out := make(chan string)
	go func() {
		for line := range in {
			// parseSomeHow
			// out <- parsedLine
		}
		close(out)
	}()
	return out
}

func saveLine(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for line := range in {
			// save line to file
		}
		close(out)
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
