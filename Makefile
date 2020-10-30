ecoli:
	go run . -encode -in=./testdata/E.coli -out=./testdata/E.coli.lzw && go run . -decode -in=./testdata/E.coli.lzw -out=./testdata/E.coli.decoded

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out
