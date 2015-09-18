package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/BatchLabs/charlatan"
	"github.com/BatchLabs/charlatan/record"
)

func main() {

	if len(os.Args) != 2 {
		usage()
		return
	}

	s := os.Args[1]

	// Creates the query

	query, err := charlatan.QueryFromString(s)
	if err != nil {
		fmt.Println(">>> ", err)
		return
	}

	reader, err := os.Open(query.From())
	if err != nil {
		fmt.Printf(">>> Error opening %s: %v\n", query.From(), err)
		return
	}

	executeRequest(csv.NewReader(reader), query)
}

func executeRequest(reader *csv.Reader, query *charlatan.Query) {

	fmt.Println("$ ", query)
	fmt.Println("$")

	var header []string
	var limit int64

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

		records, err := reader.Read()

		// end of file
		if err == io.EOF {
			return
		}
		// unknown error
		if err != nil {
			fmt.Println(err)
			return
		}

		// set the header
		if header == nil {
			header = records
			continue
		}

		// creates the record
		r := record.NewCSVRecordWithHeader(records, header)
		// evaluate the query
		match, err := query.Evaluate(r)
		if err != nil {
			fmt.Println(">>> Error while evaluating the query at line", line, ":", err)
			return
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
			fmt.Println(">>> Error while extracting the fields ", err)
			return
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

func usage() {
	fmt.Printf("Usage of %s\n", os.Args[0])
	fmt.Printf("%s query\n", os.Args[0])
	fmt.Printf("  query : the query to execute on a file (from being the file name)\n")
}
