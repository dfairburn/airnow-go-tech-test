.PHONY: doc run test clean
cov:
	go test -v -coverpkg ./... ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "main.go" > cover.out
	go tool cover -html=cover.out

clean:
	rm -f ./cover.out ./cover.out.tmp ./bin/spider

test:
	go test ./...

run:
	go run ./cmd/main.go

build:
	go build -o bin/crawler ./cmd/main.go