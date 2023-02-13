fuzz:
	@go test -fuzz FuzzSplit -fuzztime 10s
	@go test -fuzz FuzzAllocate -fuzztime 10s
	@go test -fuzz FuzzDiscount -fuzztime 10s

test:
	@go test -race -v -coverprofile cov.out -cpuprofile cpu.out -memprofile mem.out
	@go tool cover -html cov.out
