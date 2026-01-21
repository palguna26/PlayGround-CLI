# PlayGround CLI Makefile
# For local development and release building

BINARY_NAME := pg
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_FLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

# Platforms for release
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build clean test install release checksums

# Default: build for current platform
all: build

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/pg

# Install locally
install: build
	mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME) || \
	mv $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Build for all platforms
release: clean
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		CGO_ENABLED=0 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)_$${platform%/*}_$${platform#*/}/$(BINARY_NAME)$$([ $${platform%/*} = windows ] && echo .exe) ./cmd/pg; \
		echo "Built: $${platform}"; \
	done

# Create release archives
package: release
	@cd dist && for dir in */; do \
		name=$${dir%/}; \
		if echo "$$name" | grep -q windows; then \
			zip -r "$$name.zip" "$$name"; \
		else \
			tar -czvf "$$name.tar.gz" "$$name"; \
		fi; \
	done
	@echo "Packages created in dist/"

# Generate checksums
checksums: package
	@cd dist && sha256sum *.tar.gz *.zip > checksums.txt
	@echo "Checksums written to dist/checksums.txt"

# Full release build
dist: checksums
	@echo "Release artifacts ready in dist/"
	@ls -la dist/
