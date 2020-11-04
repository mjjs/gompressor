ecoli:
	go run . -encode -in=./testdata/E.coli -out=./testdata/E.coli.lzw && go run . -decode -in=./testdata/E.coli.lzw -out=./testdata/E.coli.decoded

helloworld:
	go run . -encode -in=./testdata/helloworld.txt -out=./testdata/helloworld.txt.lzw && go run . -decode -in=./testdata/helloworld.txt.lzw -out=./testdata/helloworld.txt.decoded

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

test:
	go test ./... -cover

.PHONY: ecoli coverage test
