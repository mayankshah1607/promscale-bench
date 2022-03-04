package main

import (
	"flag"

	"log"

	"github.com/mayankshah1607/promscale-bench/pkg/cli"
)

type opts struct {
	n       int
	csvDir  string
	promURL string
	stats   string
}

var cliOpts = &opts{}

func main() {
	err := cli.RunBenchmark(cliOpts.n,
		cliOpts.csvDir,
		cliOpts.promURL,
		cliOpts.stats)
	if err != nil {
		log.Fatal(err)
	}
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
