binary_name := "gitx"
cmd_path := "./cmd/gitx"
build_path := "./build"
app_version := `git describe --tags --abbrev=0 2>/dev/null || echo "dev"`
ldflags := "-X main.version=" + app_version

default: build

# Syncs dependencies
[group('core')]
sync:
    @echo "Syncing dependencies..."
    @go mod tidy
    @echo "Dependencies synced."

# Builds the binary, runs sync first
[group('core')]
build: sync
    @echo "Building the application..."
    @mkdir -p {{ build_path }}
    @go build -ldflags "{{ ldflags }}" -o {{ build_path }}/{{ binary_name }} {{ cmd_path }}
    @echo "Binary available at {{ build_path }}/{{ binary_name }}"

# Runs the application, runs build first
[group('core')]
run: build
    @echo "Running {{ binary_name }}..."
    @{{ build_path }}/{{ binary_name }}

# Installs the binary to GOPATH/bin, runs build first
[group('core')]
install: build
    @echo "Installing {{ binary_name }}..."
    @go install {{ ldflags }} {{ cmd_path }}
    @echo "{{ binary_name }} installed successfully"

# Print help message, runs build first
[group('core')]
help: build
    @{{ build_path }}/{{ binary_name }} --help

# Runs all tests
[group('dev')]
test:
    @echo "Running tests..."
    @go test -v ./...

# Runs golangci-lint
[group('dev')]
ci:
    @echo "Running golangci-lint..."
    @golangci-lint run

# Format Go code
[group('dev')]
fmt:
    @echo "Formatting Go code..."
    @go fmt ./...

# Static analysis
[group('dev')]
vet:
    @echo "Running go vet..."
    @go vet ./...

# Test for any race conditions
[group('dev')]
test-race:
    @echo "Running race tests..."
    @go test -race ./...

# Coverage summary + artifact
[group('dev')]
cover:
    @echo "Running coverage..."
    @go test -coverprofile=coverage.out ./...
    @go tool cover -func=coverage.out

# Run all checks: fmt, vet, test, ci
[group('dev')]
check: fmt vet test ci
    @echo "All checks passed."

# Print binary version, runs build first
[group('dev')]
version:
    @{{ build_path }}/{{ binary_name }} --version

# Clean and rebuild the binary
[group('maintenance')]
rebuild: clean build

# Cleans the build directory
[group('maintenance')]
clean:
    @echo "Cleaning up..."
    @rm -rf {{ build_path }}
    @echo "Cleanup complete."
