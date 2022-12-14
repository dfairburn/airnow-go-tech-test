.PHONY: doc run test clean
cov:
	go test -v -coverpkg ./... ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "main.go" > cover.out
	go tool cover -html=cover.out

clean:
	rm -f ./cover.out ./cover.out.tmp ./bin/crawler

test:
	go test ./...

build:
	go build -o bin/crawler ./cmd/main.go