fuzz:
	@go test -fuzz FuzzMoneySplit -fuzztime 10s
	@go test -fuzz FuzzMoneyAllocate -fuzztime 10s
	@go test -fuzz FuzzBigMoneySplit -fuzztime 10s
	@go test -fuzz FuzzBigMoneyAllocate -fuzztime 10s
	@go test -fuzz FuzzBigMoneyDiscount -fuzztime 10s
