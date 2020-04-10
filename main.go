package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

// convert 2d array of strings into [] of problem structs
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

// reorder problems psuedo-randomly
func shuffleProblems(p []problem) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)

	for index := 0; index < len(p); index++ {
		choice := r.Intn(len(p) - index)

		p[index], p[choice+index] = p[choice+index], p[index]
	}
}

// read contents of a csv file to generate a list of problems for the quiz
func generateProblems(csvFilename string) []problem {
	file, err := os.Open(csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Error: could not open file %s", csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}

	return parseLines(lines)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {

	// read in command line flags
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format question,answer")
	shuffle := flag.Bool("shuffle", false, "shuffle the problems to appear in random order")
	timeLimit := flag.Int("limit", 30, "time limit in seconds for each problem")
	flag.Parse()

	// generate problems from
	problems := generateProblems(*csvFilename)

	if *shuffle {
		fmt.Println(".. shuffling problems")
		shuffleProblems(problems)
	}

	correct := 0
	expired := false
	ans := make(chan string)
	for i, p := range problems {
		if expired {
			break
		}

		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ans <- answer
		}()

		select {
		case answer := <-ans:
			if answer == p.a {
				correct++
			}
		case <-time.After(time.Duration(*timeLimit) * time.Second):
			fmt.Println("Times Up!")
			expired = true
		}
	}

	fmt.Printf("You answered %d/%d correct!\n", correct, len(problems))
}
