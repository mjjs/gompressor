test:
	go test ./... -cover -coverprofile=coverage.out && go tool cover -html=coverage.out -o ./doc/coverage.html && rm coverage.out

.PHONY: ecoli coverage test
