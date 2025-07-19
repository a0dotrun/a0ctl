# Binary name
BINARY_NAME = a0

# Build the CLI
build:
	go build -o $(BINARY_NAME) ./cmd/a0ctl/main.go

# Run the CLI with optional args: `make run ARGS="help"`
a0: build
	./$(BINARY_NAME) $(ARGS)

# Clean up the binary
clean:
	rm -f $(BINARY_NAME)
