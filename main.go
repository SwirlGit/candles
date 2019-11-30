package main

func readFile(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		fileName, ok := <-in
		if ok {
			// scan file
		}
		close(out)
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
