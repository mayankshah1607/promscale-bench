# promscale-bench

A CLI tool to benchmark PromQL queries against Promscale.

## Installation

Run the following make target to build the CLI binary:
```bash
$ make build
```

This will build the CLI binary and place it in the `bin/` folder.

## Usage

```bash
$ ./bin/promscale-bench -n 3 \
--dir /path/to/csv-queries \
--prom-url <Promscale URL> \
--stats min,max,med,avg,p90,p99,errs
```

### CLI flags

| Name        | Description                                           | Default                        |
|-------------|-------------------------------------------------------|--------------------------------|
| `--n`        | Number of concurrent workers                          | 3                              |
| `--dir`      | Path the the CSV file containing queries to benchmark | ""                             |
| `--prom-url` | URL of the running Promscale / Prometheus instance    | "http://localhost:9201"               |
| `--stats`    | Comma separated list of stats to measure              | "min,max,med,avg,p90,p99,errs" |

### Run a sample benchmark

1. Clone this repo
```bash
$ git clone github.com/mayankshah1607/promscale-bench
```

2. Setup TimescaleDB and Promscale locally using Docker:
```bash
$ make dev
```

3. Run a sample bench mark using the data and queries provided in the `data/` directory:
```bash
$ make run-sample-benchmark
```

This step ingests the sample data into Promscale, builds the CLI binary and runs it against the installation of Promscale. This would give you an output similar to:

```
Processed 11 queries in 31.883343ms

Min Query time:         6.206807ms
Max Query time:         11.340485ms
Median Query time:      [7.631863ms]
Avg Query Time:         8.328757ms
90th Percentile:        11.273614ms
90th Percentile:        11.340485ms
Total Errors:           0
```

4. Clean up the environment:
```bash
$ make clean
```
