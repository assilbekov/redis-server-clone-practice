run: build
	@./bin/goredix --listenAddr :5001

build:
	@go build -o bin/goredix .
