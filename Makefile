all: build

build:
	go build -o bin/load-balancer .

clean:
	rm -f bin/load-balancer

run: build
	./bin/load-balancer

.PHONY: all build run
