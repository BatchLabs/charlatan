package main

import (
	"encoding/json"
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
	var limit int64

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

	decoder := json.NewDecoder(reader)

	hasLimit := query.HasLimit()

	if hasLimit {
		limit = query.Limit()

		if limit <= 0 {
			return
		}
	}

	line := 0
	offset := query.StartingAt()

	for {
		r, err := record.NewJSONRecordFromDecoder(decoder)

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

		offset--

		if offset > 0 {
			continue
		}

		line++

		// extract the fields
		values, err := query.FieldsValues(r)
		if err != nil {
			fatalf("Error while extracting the fields at line %d: %v\n", line, err)
		}

		fmt.Println("# ", values)
		if hasLimit {
			limit--
			if limit <= 0 {
				break
			}
		}
	}
}
