package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var csvFilename string
var timeLimit int

func init() {
	flag.StringVar(&csvFilename, "csv", "problem.csv", "a csv file in the format of 'question, answer'")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds")
}

func main() {

	flag.Parse()

	file, err := os.Open(csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provide CSV file.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.quiz)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			return
		case answer := <-answerChannel:
			if answer == p.answer {
				fmt.Println("Correct!")
				correct++
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d.\n", correct, len(lines))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			quiz:   line[0],
			answer: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	quiz   string
	answer string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
