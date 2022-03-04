NETWORK_NAME=promscale-timescale
PROM_URL=localhost:9201
BIN_DIR?=$(PWD)/bin
SAMPLE_CSV_DIR?=$(PWD)/data/obs-queries.csv

.PHONY: dev
dev:
	docker network create --driver bridge $(NETWORK_NAME);
	docker run --name timescaledb \
	-e POSTGRES_PASSWORD=admin \
	-d \
	-p 5432:5432 \
	--network $(NETWORK_NAME) \
	timescaledev/promscale-extension:latest-ts2-pg13 postgres -csynchronous_commit=off;
	@sleep 10;
	docker run --name promscale \
	-d \
	-p 9201:9201 \
	--network $(NETWORK_NAME) \
	timescale/promscale:latest \
	-db-password=admin \
	-db-port=5432 \
	-db-name=postgres \
	-db-host=timescaledb \
	-db-ssl-mode=allow;

.PHONY: clean
clean:
	docker stop timescaledb && docker rm timescaledb;
	docker stop promscale && docker rm promscale;
	docker network rm $(NETWORK_NAME);

.PHONY: ingest
ingest:
	@echo "Ingesting sample data..."
	curl -v \
	-H "Content-Type: application/x-protobuf" \
	-H "Content-Encoding: snappy" \
	-H "X-Prometheus-Remote-Write-Version: 0.1.0" \
	--request POST \
	--data-binary "@real-dataset.sz" \
	"$(PROM_URL)/write";

.PHONY: build
build:
	go build -o $(BIN_DIR)/promscale-bench cli/main.go

.PHONY: run-sample-benchmark
run-sample-banchmark:
	$(MAKE) ingest;
	$(MAKE) build;
	./bin/promscale-bench -dir $(SAMPLE_CSV_DIR)
