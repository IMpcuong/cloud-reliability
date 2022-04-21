# Global variables
GOCMD       = go
BINARY_NAME = pdpapp
EXTENSION   = exe

# List of usual golang commands
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTIDY  = $(GOCMD) mod tidy
GOTEST  = $(GOCMD) test
GOGET   = $(GOCMD) get

# Default task
all: build

# Build task
build:
	$(GOBUILD) -o $(BINARY_NAME).$(EXTENSION) -v

# Test task
test:
	$(GOTEST) -v ./...

# Target with to execute command from specific task
# Note: not really necessary (cause we have been specific the shell name)
.PHONY: clean
# Clean task
clean:
ifeq ("$(wildcard $($(BINARY_NAME).$(EXTENSION)))", "")
	$(GOCLEAN)
	$(GOTIDY)
	powershell rm $(BINARY_NAME).$(EXTENSION)
else
	$(GOCLEAN)
	$(GOTIDY)
endif

# Run task
run:
	$(GOBUILD) -o $(BINARY_NAME).$(EXTENSION) -v ./...
	./$(BINARY_NAME).$(EXTENSION)

# Dependencies
deps:
	$(GOGET) github.com/urfave/cli
	$(GOGET) golang.org/x/exp/constraints

# Update all dependencies
update:
	$(GOGET) -u
