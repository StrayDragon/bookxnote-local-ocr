.PHONY: all build prepare package clean

BUILD_OUTPUT_DIR ?= build
BINARY_SUFFIX=$(if $(filter windows,$(GOOS)),.exe)

BINARIES := bookxnote-local-ocr server certgen

ifdef GOOS
ifdef GOARCH
    PLATFORM_SUFFIX := $(GOOS)-$(GOARCH)
    BUILD_OUTPUT_DIR := $(BUILD_OUTPUT_DIR)/bookxnote-local-ocr-$(PLATFORM_SUFFIX)
    ARCHIVE_NAME := bookxnote-local-ocr-$(PLATFORM_SUFFIX)
    ARCHIVE_EXT := $(if $(filter windows,$(GOOS)),.zip,.tar.gz)
endif
endif

all: package

build:
	@echo "building..."
	@mkdir -p $(BUILD_OUTPUT_DIR)
	@for binary in $(BINARIES); do \
		echo "building $$binary..."; \
		go build -o $(BUILD_OUTPUT_DIR)/$$binary$(BINARY_SUFFIX) cmd/$$binary/main.go; \
	done

prepare: build
	@echo "prepare..."
	@cp README.md LICENSE $(BUILD_OUTPUT_DIR)/
	@mkdir -p $(BUILD_OUTPUT_DIR)/docs
	@cp docs/tutorial.md $(BUILD_OUTPUT_DIR)/docs/tutorial.md
	@cp config.yml $(BUILD_OUTPUT_DIR)/

package: prepare
	@echo "packaging..."
	@cd $(dir $(BUILD_OUTPUT_DIR)) && \
	if [ "$(GOOS)" = "windows" ]; then \
		zip -r $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(notdir $(BUILD_OUTPUT_DIR)); \
	else \
		tar czf $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(notdir $(BUILD_OUTPUT_DIR)); \
	fi
	@echo "done: $(ARCHIVE_NAME)$(ARCHIVE_EXT)"

clean:
	@echo "clean..."
	@rm -rf $(BUILD_OUTPUT_DIR)
