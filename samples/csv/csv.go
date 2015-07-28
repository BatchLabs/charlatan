package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/batchlabs/charlatan"
	"github.com/batchlabs/charlatan/record"
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

	var header *[]string
	line := 0

	for {

		records, err := reader.Read()
		line++

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
			header = &records
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

		// extract the fields
		values, err := query.FieldsValues(r)
		if err != nil {
			fmt.Println(">>> Error while extracting the fields ", err)
			return
		}

		fmt.Println("# ", values)
	}
}

func usage() {
	fmt.Printf("Usage of %s\n", os.Args[0])
	fmt.Printf("%s query\n", os.Args[0])
	fmt.Printf("  query : the query to execute on a file (from being the file name)\n")
}
