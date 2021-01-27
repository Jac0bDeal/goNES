all: clean lint test goNES

clean:
	@echo "Cleaning bin/..."
	@rm -rf bin/*

dependencies:
	@echo "Installing project dependencies..."
	@go get -u golang.org/x/lint/golint

goNES:
	@echo "Building goNES binary for use on local system..."
	@go build -o bin/goNES ./cmd/goNES

lint:
	@echo "Running linters..."
	@go vet ./...
	@golint ./...

test:
	@echo "Running tests..."
	@go test ./...
