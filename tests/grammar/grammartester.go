package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/batchlabs/charlatan"
)

const testsCount = 10

func generateQueryString() (string, error) {
	c := exec.Command("abnfgen", "-l", "-s", "query", "grammar.abnf")
	c.Stderr = nil

	out, err := c.Output()
	return string(out), err
}

func runTest() (string, error) {
	s, err := generateQueryString()
	if err != nil {
		log.Fatalf("Query generation error: %v", err)
	}

	_, err = charlatan.QueryFromString(s)
	return s, err
}

func main() {
	var err error
	var good, count, tc int64

	if len(os.Args) == 2 {
		if count, err = strconv.ParseInt(os.Args[1], 10, 64); err != nil {
			count = testsCount
		}
	} else {
		count = testsCount
	}

	for tc = 0; tc < count; tc++ {
		s, err := runTest()
		if err == nil {
			good++
			log.Printf("[OK] %s", s)
		} else {
			log.Printf("[KO] %s", s)
			log.Printf("     ERR: %v", err)
		}
	}

	log.Println("")
	log.Printf("Results: %d/%d", good, count)
}
