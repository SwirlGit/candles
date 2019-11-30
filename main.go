package main

func readFile(in <-chan string) <-chan string {
	out := make(chan string)
	return out
}

func parseLine(in <-chan string) <-chan string {
	out := make(chan string)
	return out
}

func saveLine(in <-chan string) <-chan string {
	out := make(chan string)
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
