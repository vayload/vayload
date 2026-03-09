.PHONY: all bench bench-cold bench-warm bench-mem bench-concurrent bench-pool bench-sustained bench-all tidy

tidy:
	go mod tidy

# Run ALL benchmarks + tests (recommended first run)
all: tidy bench-all

# Full suite: benchmarks + human-readable tests
bench-all:
	@echo "=== Running all benchmarks + detailed tests ==="
	go test ./performance/ \
		-v \
		-run "TestColdStartDetailed|TestSustainedLoad" \
		-bench "." \
		-benchmem \
		-benchtime=3s \
		-count=1 \
		-cpu=1,2,4,8 \
		| tee results.txt
	@echo "\n✅ Results saved to results.txt"

# Cold start only
bench-cold:
	go test ./performance/ \
		-bench "BenchmarkColdStart" \
		-benchmem \
		-benchtime=3s \
		-count=3

# Warm execution only
bench-warm:
	go test ./performance/ \
		-bench "BenchmarkWarm" \
		-benchmem \
		-benchtime=5s \
		-count=3

# Memory pressure
bench-mem:
	go test ./performance/ \
		-bench "BenchmarkMemory" \
		-benchmem \
		-benchtime=3s \
		-count=3

# Concurrent
bench-concurrent:
	go test ./performance/ \
		-bench "BenchmarkConcurrent" \
		-benchmem \
		-benchtime=5s \
		-count=3 \
		-cpu=1,2,4,8

# Pool / reuse pattern
bench-pool:
	go test ./performance/ \
		-bench "BenchmarkPool" \
		-benchmem \
		-benchtime=5s \
		-count=3

# Human-readable sustained load + cold start report
bench-sustained:
	go test ./performance/ \
		-v \
		-run "TestSustainedLoad|TestColdStartDetailed" \
		-timeout 60s

# Profile CPU (outputs cpu.prof)
profile-cpu:
	go test ./performance/ \
		-bench "BenchmarkWarm_GopherLua_Fib|BenchmarkWarm_Goja_Fib|BenchmarkWarm_Tengo_Fib" \
		-benchmem \
		-benchtime=5s \
		-cpuprofile=cpu.prof
	@echo "Run: go tool pprof cpu.prof"

# Profile memory (outputs mem.prof)
profile-mem:
	go test ./performance/ \
		-bench "BenchmarkColdStart" \
		-benchmem \
		-benchtime=5s \
		-memprofile=mem.prof
	@echo "Run: go tool pprof mem.prof"

# Profile with visualizer (needs graphviz)
profile-viz: profile-cpu
	go tool pprof -http=:6060 cpu.prof
