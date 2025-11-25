.PHONY: run build clean

# Default target
all: run

# Run the application
run:
	go run ./cmd/watchtower

# Build the binary
build:
	go build -o pr-watchtower ./cmd/watchtower

# Clean build artifacts
clean:
	rm -f pr-watchtower
	rm -f *.db
