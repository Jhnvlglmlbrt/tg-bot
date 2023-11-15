build:
	@go build -o bin/app

run: build
	@./bin/app -h "$(host)" -t "$(token)"

.PHONY: build run
