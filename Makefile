.PHONY: build

BIN_DIR := $(shell pwd)/bin
GO_BUILD_TAGS := -tags origin
GOOS := $(shell go env GOOS)
TARGET := github.com/dxfeed/dxfeed-graal-go-api/cmd/tools

define GO_BUILD_CMD
	@echo "Building for $(1)..."
	go build $(GO_BUILD_TAGS) -v ./...
	GOBIN=$(BIN_DIR) go install $(GO_BUILD_TAGS) $(TARGET)
endef

define GO_BUILD_CMD_MACOS
	@echo "Building for $(1)..."
	CGO_LDFLAGS_ALLOW=-Wl,-rpath,@executable_path/ go build $(GO_BUILD_TAGS) -v ./...
	GOBIN=$(BIN_DIR) CGO_LDFLAGS_ALLOW=-Wl,-rpath,@executable_path/ go install $(GO_BUILD_TAGS) $(TARGET)
	codesign -f -s - $(BIN_DIR)/*
endef

# Main build target
build:
ifeq ($(GOOS),windows)
	$(call GO_BUILD_CMD,windows)
else ifeq ($(GOOS),linux)
	$(call GO_BUILD_CMD,linux)
else ifeq ($(GOOS),darwin)
	$(call GO_BUILD_CMD_MACOS,macOS)
endif
	cp internal/native/graal/*DxFeedGraalNativeSdk.* $(BIN_DIR)/
