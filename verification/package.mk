.PHONY: benchmark concurrency docs package-specific reproducible

package-specific: concurrency reproducible

concurrency:
	go test -race ./... \
		-run 'Test(Progress|DynamicOptions|StatusStream|TaskGroup)' -count=20

docs:
	./scripts/check-docs.sh

reproducible:
	./scripts/check-reproducible-source.sh

benchmark:
	go test ./... -run '^$$' -bench Benchmark -benchmem -benchtime=100ms
	cd benchmarks/comparison && go mod tidy -diff
	cd benchmarks/comparison && go test ./... -run '^$$'
	cd benchmarks/comparison && ./run-benchmarks.sh
	cd benchmarks/comparison && ./measure-binaries.sh
