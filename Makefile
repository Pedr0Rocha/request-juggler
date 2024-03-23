all: build

build:
	go build -o bin/request-juggler .

clean:
	rm -f bin/request-juggler

run: build
	./bin/request-juggler

.PHONY: all build run
