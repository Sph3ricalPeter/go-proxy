all: clean test build run

clean:
	rm -rf bin

test:
	go test

build:
	go build -o bin/main main.go

run:
	./bin/main

.PHONY: all build run