SHELL := /bin/bash

# Run individual services with their own env files
run-user:
	cd user-management && go mod tidy && go run main.go
.PHONY: run-user

# Run all services (in separate terminals is ideal)
run-all:
	@echo "Run each service in a new terminal for proper logging."
	@echo "Use: make run-user | make run-auth | make run-payment"
.PHONY: run-all

# Clean binaries (if you build any)
clean:
	find . -type f -name "*.out" -delete
.PHONY: clean
