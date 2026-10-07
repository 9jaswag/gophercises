// Quiz game

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Quiz Game")
	// read quiz from file, line by line
	// keep track of how many questions they ger right, wrong
	// ask next question immediately
	// output total correct total questions

	ra, tq := 0, 0

	bs, err := os.ReadFile("1/problems.csv")
	if err != nil {
		log.Println("Error:", err)
	}
	data := string(bs)
	r := csv.NewReader(strings.NewReader(data))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		tq++

		var ans string
		q, a := record[0], record[1]

		fmt.Println(q, "=")
		fmt.Scan(&ans)

		if strings.TrimSpace(ans) == a {
			ra++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", ra, tq)
}
