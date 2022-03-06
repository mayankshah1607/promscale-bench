package main

import (
	"flag"
	"fmt"

	"log"

	"github.com/mayankshah1607/promscale-bench/pkg/cli"
	prom "github.com/mayankshah1607/promscale-bench/pkg/prometheus"
)

type opts struct {
	n       int
	csvDir  string
	promURL string
	stats   string
}

var cliOpts = &opts{}

func main() {

	// validate CLI args
	if err := validateArgs(cliOpts); err != nil {
		log.Fatal(err)
	}

	// Test the connection to Prometheus
	if err := prom.ConnectionTest(cliOpts.promURL); err != nil {
		log.Fatal(err)
	}

	// Run the benchmarks
	if err := cli.RunBenchmark(cliOpts.n,
		cliOpts.csvDir,
		cliOpts.promURL,
		cliOpts.stats); err != nil {
		log.Fatal(err)
	}
}

func validateArgs(opts *opts) error {
	if opts.csvDir == "" {
		return fmt.Errorf("`-dir` cannot be empty. Please specify a CSV file")
	}

	if opts.promURL == "" {
		return fmt.Errorf("`prom-url` cannot be empty." +
			"Please specify the URL to your Prometheus instance")
	}
	return nil
}

func init() {
	flag.IntVar(&cliOpts.n, "n", 3,
		"Number of concurrent workers")
	flag.StringVar(&cliOpts.csvDir, "dir", "",
		"Diretory of the CSV file with PromQL queries")
	flag.StringVar(&cliOpts.promURL, "prom-url", "http://localhost:9201",
		"URL of Prometheus/Promscale instance")
	flag.StringVar(&cliOpts.stats, "stats", "min,max,med,avg,p90,p99,errs",
		"Comma-separated list of stats to measure. Available: num,total_time,errs,min,max,avg,p99,p90")

	flag.Parse()
}
