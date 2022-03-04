package cli

import (
	"fmt"
	"math"
	"time"

	prom "github.com/mayankshah1607/promscale-bench/pkg/prometheus"
)

type statComputeFunc func(sortedResults []prom.PromQLResponse)

func min(sortedResults []prom.PromQLResponse) {
	fmt.Printf("Min Query time:\t\t%s\n", sortedResults[0].ExecTime)
}

func max(sortedResults []prom.PromQLResponse) {
	l := len(sortedResults)
	fmt.Printf("Max Query time:\t\t%s\n", sortedResults[l-1].ExecTime)
}

func median(sortedResults []prom.PromQLResponse) {
	l := len(sortedResults)
	medianQueryTime := []time.Duration{
		sortedResults[l/2].ExecTime}
	if l%2 == 0 {
		medianQueryTime = append(medianQueryTime,
			sortedResults[(l/2)+1].ExecTime)
	}
	fmt.Printf("Median Query time:\t%v\n", medianQueryTime)
}

func p90(sortedResults []prom.PromQLResponse) {
	l := len(sortedResults)
	p90Index := int(math.Ceil(0.9*float64(l))) - 1
	fmt.Printf("90th Percentile:\t%s\n", sortedResults[p90Index].ExecTime)
}

func p99(sortedResults []prom.PromQLResponse) {
	l := len(sortedResults)
	p90Index := int(math.Ceil(0.99*float64(l))) - 1
	fmt.Printf("90th Percentile:\t%s\n", sortedResults[p90Index].ExecTime)
}

func errorCount(sortedResults []prom.PromQLResponse) {
	errCount := 0
	for _, r := range sortedResults {
		if r.Response.StatusCode != 200 || r.Err != nil {
			errCount += 1
		}
	}
	fmt.Printf("Total Errors:\t\t%d\n", errCount)
}

func avg(sortedResults []prom.PromQLResponse) {
	l := len(sortedResults)
	sum := 0.0
	for _, t := range sortedResults {
		sum = sum + float64(t.ExecTime)
	}
	avg := sum / float64(l)
	fmt.Printf("Avg Query Time:\t\t%s\n", time.Duration(avg))
}

func getStatComputeFuncs(stats []string) []statComputeFunc {
	computeFuncs := []statComputeFunc{}
	for _, stat := range stats {
		switch stat {
		case Min:
			computeFuncs = append(computeFuncs, min)
		case Max:
			computeFuncs = append(computeFuncs, max)
		case Median:
			computeFuncs = append(computeFuncs, median)
		case P90:
			computeFuncs = append(computeFuncs, p90)
		case P99:
			computeFuncs = append(computeFuncs, p99)
		case ErrCount:
			computeFuncs = append(computeFuncs, errorCount)
		case Avg:
			computeFuncs = append(computeFuncs, avg)
		}
	}
	return computeFuncs
}
