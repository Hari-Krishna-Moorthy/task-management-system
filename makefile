# Variables
APP_NAME := task-management-system
MAIN_DB_URI := $(shell grep MAIN_DB_URI .env | sed 's/MAIN_DB_URI=//')
TEST_DB_URI := $(shell grep TEST_DB_URI .env | sed 's/TEST_DB_URI=//')

# Go build and test commands
GO_CMD := go
DOCKER_COMPOSE_CMD := docker-compose
DOCKER_CMD := docker


# Set up the project (install dependencies, setup database)
setup: 
	@echo "Setting up the project..."
	@$(GO_CMD) mod tidy


# Run tests using the test database configuration
test: 
	@echo "Running tests with test database..."
	@MAIN_DB_URI=$(TEST_DB_URI) $(GO_CMD) test ./internal/app/test/... -v

# Run tests and generate a coverage report
test-report:
	@echo "Running tests and generating coverage report..."
	@MAIN_DB_URI=$(TEST_DB_URI) $(GO_CMD) test ./internal/app/test/... -coverprofile=coverage.out
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run the application locally
run:
	@echo "Running the application locally..."
	@$(GO_CMD) run cmd/server/main.go

# Build and run the application with Docker
docker-build:
	@echo "Building Docker image..."
	@$(DOCKER_CMD) build -t $(APP_NAME) .

docker-run: docker-build
	@echo "Running Docker container..."
	@$(DOCKER_CMD) run --env-file .env -p 3000:3000 $(APP_NAME)

# Clean up generated files
clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.html


# Display lint errors
lint:
	@echo "Linting..."
	golangci-lint run  

# Display help
help:
	@echo "Makefile commands:"
	@echo "  setup            - Set up the project (dependencies, database)"
	@echo "  test             - Run tests using the test database"
	@echo "  test-report      - Run tests and generate a coverage report"
	@echo "  run              - Run the application locally"
	@echo "  docker-build     - Build the Docker image"
	@echo "  docker-run       - Run the application in a Docker container"
	@echo "  clean            - Clean up generated files"

