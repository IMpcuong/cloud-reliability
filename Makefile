# Global variables
GOCMD = go
BINARY_NAME = pdp_bc.out

# List of usual golang commands
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTIDY = $(GOCMD) mod tidy
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Default task
all: build

# Build task
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Test task
test:
	$(GOTEST) -v ./...

# Target with to execute command from specific task
# Note: not really necessary (cause we have been specific the shell name)
.PHONY: clean
# Clean task
clean:
	$(GOCLEAN)
	$(GOTIDY)
	powershell rm *.out

# Run task
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Dependencies
deps:
	$(GOGET) github.com/urfave/cli

# Update all dependencies
update:
	$(GOGET) -u
