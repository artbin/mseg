SHELL := /bin/zsh
GO    ?= go

.PHONY: help all test bench bench1 plots plot-deps fmt vet tidy visualize clean

help:
	@echo "Targets:"
	@echo "  test       - run unit tests"
	@echo "  bench      - run benchmarks (count=5) and write bench/results.csv"
	@echo "  bench1     - quick single-run benchmarks (count=1)"
	@echo "  plots      - generate plots from bench/results.csv into bench/plots"
	@echo "  plot-deps  - install Python plotting deps via uv"
	@echo "  fmt vet    - format and vet code"
	@echo "  tidy       - go mod tidy"
	@echo "  visualize  - generate memviz PNGs via cmd/visualize"
	@echo "  clean      - remove bench artifacts"

test:
	cd segment && $(GO) test ./...
	cd mlist && $(GO) test ./...
	cd marray && $(GO) test ./...
	cd pflist && $(GO) test ./...
	cd pfqueue && $(GO) test ./...
	cd pfdeque && $(GO) test ./...

bench:
	cd cmd/benchrun && $(GO) run . -count 5 -out ../../bench/results.csv -pkg ../..

bench1:
	cd cmd/benchrun && $(GO) run . -count 1 -out ../../bench/results.csv -pkg ../..

plots:
	uv run --project bench bench/plot.py --in bench/results.csv --out bench/plots

plot-deps:
	uv sync --project bench

fmt:
	cd segment && $(GO) fmt ./...
	cd mlist && $(GO) fmt ./...
	cd marray && $(GO) fmt ./...
	cd pflist && $(GO) fmt ./...
	cd pfqueue && $(GO) fmt ./...
	cd pfdeque && $(GO) fmt ./...
	cd cmd/visualize && $(GO) fmt ./...

vet:
	cd segment && $(GO) vet ./...
	cd mlist && $(GO) vet ./...
	cd marray && $(GO) vet ./...
	cd pflist && $(GO) vet ./...
	cd pfqueue && $(GO) vet ./...
	cd pfdeque && $(GO) vet ./...
	cd cmd/visualize && $(GO) vet ./...

tidy:
	cd segment && $(GO) mod tidy
	cd mlist && $(GO) mod tidy
	cd marray && $(GO) mod tidy
	cd pflist && $(GO) mod tidy
	cd pfqueue && $(GO) mod tidy
	cd pfdeque && $(GO) mod tidy
	cd cmd/visualize && $(GO) mod tidy

visualize:
	cd cmd/visualize && $(GO) run . -out ../../visualize

clean:
	rm -rf bench/results*.csv bench/plots
