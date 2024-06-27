package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Problem struct {
	q, a string
}

func parseProblems(records [][]string) []Problem {
	ret := make([]Problem, len(records))
	for i, r := range records {
		ret[i] = Problem{q: r[0], a: r[1]}
	}
	return ret
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "")
	timeLimit := flag.Int("limit", 30, "")
	flag.Parse()

	handleErr := func(err error) {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file, err := os.Open(*csvFileName)
	if err != nil {
		handleErr(err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		handleErr(err)
	}

	problems := parseProblems(records)

	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

	var correction int
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s\n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				handleErr(err)
			}
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("you scored %d out of %d\n", correction, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correction++
			}
		}
	}
	fmt.Printf("you scored %d out of %d\n", correction, len(problems))
}
