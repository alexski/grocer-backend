.PHONY: build

build:
	go build -o ./bin/grocer-api ./cmd/.

.PHONY: run

run: build
	./bin/grocer-api