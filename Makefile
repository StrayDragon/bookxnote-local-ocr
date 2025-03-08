.PHONY: all build prepare package clean

# Get GOOS and GOARCH from go env
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP_NAME := bookxnote-local-ocr
BINARIES := $(APP_NAME) server certgen

# Set platform suffix and build directory
PLATFORM_SUFFIX := $(GOOS)-$(GOARCH)
BUILD_ROOT := build
BUILD_OUTPUT_DIR := $(BUILD_ROOT)/$(APP_NAME)-$(PLATFORM_SUFFIX)
ARCHIVE_NAME := $(APP_NAME)-$(PLATFORM_SUFFIX)
ARCHIVE_EXT := $(if $(filter windows,$(GOOS)),.zip,.tar.gz)
BINARY_SUFFIX := $(if $(filter windows,$(GOOS)),.exe,)

all: package

evalate-linux-build-privilege:
	sudo ./$(BUILD_OUTPUT_DIR)/setup-privileges.sh

generate-gui-resources:
	go generate ./cmd/bookxnote-local-ocr/

generate-swagger-doc:
	swag init -g ./cmd/server/main.go -o internal/swagger-doc

generate: generate-gui-resources generate-swagger-doc
	@echo "generated"

build: generate clean
	@echo "building..."
	@mkdir -p $(BUILD_OUTPUT_DIR)
	@for binary in $(BINARIES); do \
		echo "building $$binary..."; \
		go build -o $(BUILD_OUTPUT_DIR)/$$binary$(BINARY_SUFFIX) ./cmd/$$binary/...; \
	done
	@cp artifact/config.yml $(BUILD_OUTPUT_DIR)/

prepare: build
	@echo "prepare..."
	@cp README.md LICENSE $(BUILD_OUTPUT_DIR)/
	@mkdir -p $(BUILD_OUTPUT_DIR)/docs
	@cp docs/tutorial.md $(BUILD_OUTPUT_DIR)/docs/tutorial.md
	@if [ "$(GOOS)" = "windows" ]; then \
		cp artifact/windows/*.bat $(BUILD_OUTPUT_DIR)/; \
	elif [ "$(GOOS)" = "linux" ]; then \
		cp artifact/linux/setup-privileges.sh $(BUILD_OUTPUT_DIR)/; \
		chmod +x $(BUILD_OUTPUT_DIR)/setup-privileges.sh; \
	fi

package: prepare
	@echo "packaging..."
	@cd $(BUILD_ROOT) && \
	if [ "$(GOOS)" = "windows" ]; then \
		zip -r $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(APP_NAME)-$(PLATFORM_SUFFIX); \
	else \
		tar czf $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(APP_NAME)-$(PLATFORM_SUFFIX); \
	fi
	@echo "done: $(ARCHIVE_NAME)$(ARCHIVE_EXT)"

dev-build: prepare evalate-linux-build-privilege
	@echo "Done"

dev-run-gui: dev-build
	./build/bookxnote-local-ocr-linux-amd64/bookxnote-local-ocr gui

clean:
	@echo "clean..."
	@rm -rf $(BUILD_ROOT)
