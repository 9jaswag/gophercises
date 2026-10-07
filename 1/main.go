// Quiz game

package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Quiz Game")

	csvFilename := flag.String("csv", "problems.csv", "a csv with the format 'problem,answer'")
	duration := flag.Int("time", 10, "quiz time limit in seconds")
	flag.Parse()

	data := parseData(*csvFilename)
	problems := parseProblems(data)
	ra := serveQuestions(problems, *duration)

	fmt.Printf("You scored %d out of %d.\n", ra, len(data))
}

type problem struct {
	q string
	a string
}

func parseData(fn string) [][]string {
	f, err := os.Open("1/" + fn)
	if err != nil {
		fmt.Printf("Unable to load csv file - %v \n", fn)
		os.Exit(1)
	}
	defer f.Close()

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		fmt.Println("Unable to parse csv")
		os.Exit(1)
	}

	return data
}

func parseProblems(data [][]string) []problem {
	problems := make([]problem, len(data))
	for i, v := range data {
		problems[i] = problem{
			q: strings.TrimSpace(v[0]),
			a: strings.TrimSpace(v[1]),
		}
	}

	return problems
}

func serveQuestions(probs []problem, d int) int {
	duration := time.Duration(d) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	ra := 0

	for i, p := range probs {
		if ctx.Err() != nil {
			break
		}

		var ans string
		fmt.Printf("Question #%v: %v=\n", i, p.q)
		fmt.Scan(&ans)

		if strings.TrimSpace(ans) == p.a {
			ra++
		}
	}

	return ra
}
