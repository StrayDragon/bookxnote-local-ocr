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

OPENAPI_OUTPUT_DIR:=internal/client/openapi
OPENAPI_CLEANUP_FILES:=docs test .openapi-generator .gitignore .travis.yml git_push.sh .openapi-generator-ignore api

all: package

evalate-linux-build-privilege:
	sudo ./$(BUILD_OUTPUT_DIR)/setup-privileges.sh

generate-gui-resources:
	go generate ./cmd/bookxnote-local-ocr/

generate-api-doc:
	mkdir -p internal/swagger-doc/openapi
	# replace to swag2op get openapi 3.0 spec for better generator experience
	# go run github.com/swaggo/swag/cmd/swag@latest init -g ./cmd/server/main.go -o internal/swagger-doc
	go run github.com/zxmfke/swagger2openapi3/cmd/swag2op@latest init -g ./cmd/server/main.go \
		-o internal/swagger-doc/ \
		-openo internal/swagger-doc/openapi

generate-api-client: generate-api-doc
	mkdir -p $(OPENAPI_OUTPUT_DIR)
	uvx openapi-generator-cli generate -i internal/swagger-doc/openapi/swagger.yaml \
		-g go -o $(OPENAPI_OUTPUT_DIR) \
		--skip-validate-spec \
		--additional-properties=disallowAdditionalPropertiesIfNotPresent=false,enumClassPrefix=true,structPrefix=true,generateInterfaces=true,isGoSubmodule=true,withGoMod=false
	rm -rf $(addprefix $(OPENAPI_OUTPUT_DIR)/,$(OPENAPI_CLEANUP_FILES))

generate: generate-gui-resources generate-api-doc generate-api-client
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
