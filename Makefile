.PHONY: build run test clean tidy

APP_NAME=learn-fiber-auth-jwt
MAIN_FILE=cmd/api/main.go

# Build the application
build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) $(MAIN_FILE)

# Run the application directly
run:
	@echo "Running..."
	@go run $(MAIN_FILE)

# Run tests
test:
	@echo "Testing..."
	@go test -v ./...

# Clean the build directory
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf bin/

# Tidy module dependencies
tidy:
	@echo "Tidying..."
	@go mod tidy