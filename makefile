build:
	@go build -o bin/gotodoapp

run: 
	@./bin/gotodoapp

test:
	@go test -v ./...