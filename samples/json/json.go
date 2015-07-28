package main

import (
	"fmt"
	"io"
	"os"

	"github.com/BatchLabs/charlatan"
	"github.com/BatchLabs/charlatan/record"
)

func usage() {
	fmt.Printf("Usage:\n\n\t%s <query>\n", os.Args[0])
	os.Exit(1)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	query, err := charlatan.QueryFromString(os.Args[1])
	if err != nil {
		fatalf("Error: %v\n", err)
	}

	reader, err := os.Open(query.From())
	if err != nil {
		fatalf("Error: %v\n", err)
	}

	defer reader.Close()

	executeRequest(reader, query)
}

func executeRequest(reader io.Reader, query *charlatan.Query) {
	line := 0

	for {
		line++

		r, err := record.NewJSONRecordFromReader(reader)

		// end of file
		if err == io.EOF {
			return
		}

		// unknown error
		if err != nil {
			fatalf("Error at line %d: %v\n", line, err)
		}

		match, err := query.Evaluate(r)
		if err != nil {
			fatalf("Error while evaluating the query at line %d: %v\n", line, err)
		}

		if !match {
			continue
		}

		// extract the fields
		values, err := query.FieldsValues(r)
		if err != nil {
			fatalf("Error while extracting the fields at line %d: %v\n", line, err)
		}

		fmt.Println("# ", values)
	}
}
