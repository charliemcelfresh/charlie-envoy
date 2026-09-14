.PHONY: test coverage coverage-html

help:
	@echo "Available targets:"
	@echo "  test           Run all Go tests"
	@echo "  coverage       Run tests and generate a coverage profile"
	@echo "  coverage-html  Run tests and generate an HTML coverage report"

# Run all Go tests
test:
	go test ./...

# Run tests and generate a coverage profile
coverage:
	go test ./... -coverprofile=coverage.out

# Run tests, generate coverage, and create an HTML report
coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html