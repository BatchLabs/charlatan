package main

import (
	"fmt"
	"os"

	"github.com/BatchLabs/charlatan"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n\t\t%s <query>", os.Args[0])
		os.Exit(1)
	}

	l := charlatan.LexerFromString(os.Args[1])

	for {
		t, err := l.NextToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v", err)
			os.Exit(1)
		}
		fmt.Printf("%s ", t)
		if t.IsEnd() {
			break
		}
	}
	fmt.Println()
}
