package cli

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	prom "github.com/mayankshah1607/promscale-bench/pkg/prometheus"
	"github.com/mayankshah1607/promscale-bench/pkg/utils"
)

// Available stats this tool can measure
var (
	Avg      = "avg"
	Median   = "med"
	Min      = "min"
	Max      = "max"
	P90      = "p90"
	P99      = "p99"
	ErrCount = "errs"
)

func RunBenchmark(workerN int, csvDir, promURL, stats string) error {

	// Parse the CSV and load the queries into a `queries` slice buffer
	queries, err := utils.ParseCSV(csvDir)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	var (
		// Total number of queries / work processed
		queueSize = len(queries)

		// Spawn new worker pool with `workerN` routines
		// and a queue size of `queueSize`
		wp = utils.NewWorkerPool(workerN, queueSize)

		// workers will push results to this channel
		resultC = make(chan prom.PromQLResponse, queueSize)

		// HTTP client used for making PromQL queries
		httpClient = http.Client{}
	)

	for _, query := range queries {
		// Create new HTTP request for each query
		req, err := prom.CreateHTTPRequest(promURL, query)
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %v", err)
		}

		// Queue the request to the task queue
		wp.Add(func() {
			resultC <- prom.ExecRequestWithClient(
				req,
				&httpClient,
			)
		})
	}
	// Close task channel once all work is loaded.
	// The task queue must only be read from now on.
	wp.CloseTaskC()

	start := time.Now()
	wp.Start() // start workers
	wp.Wait()  // wait for workers to finish

	elapsed := time.Since(start) // mark the completion time for the workers

	// Closing the results channel here is safe because
	// all workers (producers) have exited after `wp.Wait()`.
	close(resultC)

	processAndPrintBenchmarkData(queueSize, elapsed,
		strings.Split(stats, ","), resultC)
	return nil
}

func processAndPrintBenchmarkData(nQueries int,
	totalTime time.Duration, stats []string,
	results <-chan prom.PromQLResponse) {

	// load all results from the channel into a buffer slice
	resultsBuffer := []prom.PromQLResponse{}
	for result := range results {
		resultsBuffer = append(resultsBuffer, result)
	}

	// Sorting will help us compute
	// min, max, median and percentile query time
	sortedResults := resultsBuffer
	sort.Slice(sortedResults, func(i, j int) bool {
		return int(sortedResults[i].ExecTime) < int(sortedResults[j].ExecTime)
	})

	fmt.Printf("\nProcessed %d queries in %s\n\n", nQueries, totalTime)

	computeFuncs := getStatComputeFuncs(stats)
	for _, computeFunc := range computeFuncs {
		computeFunc(sortedResults)
	}
}
