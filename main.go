package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Params defines accepted fields from command line
type Params struct {
	In        *string
	Out       *string
	Delimiter *string
	Fields    *int
}

// JSONStructure defines a structure for the generated JSON
type JSONStructure struct {
	Columns []string `json:"columns"`
	Lines   []Line   `json:"lines"`
}

// Line represents a single line of the CSV.
// A line is made of one or more records
type Line struct {
	Records map[string]interface{}
}

func main() {
	fmt.Println("### CSV TO JSON ###")

	// Parse parameters from command line
	params := Params{
		In:  flag.String("i", "", "input file"),
		Out: flag.String("o", "output.json", "output file"),
	}

	flag.Parse()

	// Read the specified file
	fmt.Println("Opening file ", *params.In)
	data, err := ioutil.ReadFile(*params.In)
	check(err)

	// Read the CSV
	fmt.Println("Reading CSV")
	r := csv.NewReader(strings.NewReader(string(data)))

	// Now, we create the JSON Structure
	counter := 0
	var columns []string
	var lines []Line

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if counter == 0 {
			columns = record
		} else {
			recMap := make(map[string]interface{})
			for i, k := range columns {
				recMap[k] = record[i]
			}

			lines = append(lines, Line{Records: recMap})
		}
		counter++
	}

	// Save the JSON file
	fmt.Println("Creating the JSON file at " + *params.Out)
	jsonStr := JSONStructure{
		Columns: columns,
		Lines:   lines,
	}

	finalJSON, marsherr := json.MarshalIndent(jsonStr, "", "\t")
	check(marsherr)

	ioerr := ioutil.WriteFile(*params.Out, finalJSON, 0644)
	check(ioerr)

	fmt.Println("JSON file created.")
}
