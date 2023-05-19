package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// btw it's not recommanded to define globaly these vars
// but just for the sake of this exercise will do that xD
var (
	correctAnswers int
	totalQuestions int
)

const defualtProblemFilename = "problem.csv"

func main() {

	var (
		flagTimer            = flag.Duration("t", 30*time.Second, "The max time for quizz")
		flagProblemFilername = flag.String("p", defualtProblemFilename, "the path to problems csv file")
	)

	flag.Parse()

	fmt.Printf("the exam start with %q time duration is %v Hit enter to start the exam\n", *flagProblemFilername, *flagTimer)
	fmt.Scanln()

	f, err := os.Open(*flagProblemFilername)

	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)

	questions, err := r.ReadAll()
	totalQuestions = len(questions)
	if err != nil {
		log.Fatal(err)
	}

	quizDone := startQuizz(questions)

	quizTimer := time.NewTimer(*flagTimer).C

	select {
	case <-quizDone:
	case <-quizTimer:
	}

	fmt.Printf("result: %d / %d\n", correctAnswers, totalQuestions)

}

func startQuizz(questions [][]string) chan bool {

	done := make(chan bool)

	go func() {
		for i, record := range questions {

			question, correctAnswer := record[0], record[1]

			fmt.Printf("%d. %s = ?\n", i+1, question)
			var answer string

			if _, err := fmt.Scan(&answer); err != nil {
				fmt.Printf("Failed to scan: %v\n", err)
				return
			}
			if answer == correctAnswer {
				correctAnswers++
			}
		}

		done <- true
	}()

	return done
}
