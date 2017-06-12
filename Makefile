default: test build

build: *.go
	CGO_ENABLED=0 go build linkhide.go
	docker build -t linkhide .

test:
	go test ./...

.PHONY: clean
clean:
	rm -f linkhide
