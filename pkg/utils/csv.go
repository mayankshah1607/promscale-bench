package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	prom "github.com/mayankshah1607/promscale-bench/pkg/prometheus"
)

// ParseCSV is a helper function to parse a CSV from a given directory
// into []PromQLRequest
func ParseCSV(csvDir string) ([]prom.PromQLRequest, error) {
	var queries []prom.PromQLRequest

	csvFile, err := os.Open(csvDir)
	if err != nil {
		return queries, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '|'
	csvReader.LazyQuotes = true

	csvLines, err := csvReader.ReadAll()
	if err != nil {
		return queries, err
	}

	for _, line := range csvLines {
		if len(line) != 4 {
			return queries,
				fmt.Errorf("invalid CSV: expected 4 values, but got %d", len(line))
		}
		rec := prom.PromQLRequest{
			Query: line[0],
			Start: line[1],
			End:   line[2],
			Step:  line[3],
		}
		queries = append(queries, rec)
	}
	return queries, nil
}
