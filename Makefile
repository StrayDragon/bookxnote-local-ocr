.PHONY: all build prepare package clean

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

CLI_NAME := bookxnote-local-ocr
GUI_NAME := bookxnote-local-ocr-gui
SERVER_NAME := server
CERTGEN_NAME := certgen
APP_NAME := bookxnote-local-ocr

# Set platform suffix and build directory
PLATFORM_SUFFIX := $(GOOS)-$(GOARCH)
BUILD_ROOT := build
BUILD_OUTPUT_DIR := $(BUILD_ROOT)/$(APP_NAME)-$(PLATFORM_SUFFIX)
ARCHIVE_NAME := $(APP_NAME)-$(PLATFORM_SUFFIX)
ARCHIVE_EXT := $(if $(filter windows,$(GOOS)),.zip,.tar.gz)
BINARY_SUFFIX := $(if $(filter windows,$(GOOS)),.exe,)

API_CLIENT_NAME:=inner_server
API_CLIENT_OUTPUT_DIR:=internal/client/$(API_CLIENT_NAME)
OPENAPI_CLIENT_NAME:=openapi
OPENAPI_CLIENT_OUTPUT_DIR:=internal/client/$(OPENAPI_CLIENT_NAME)
OPENAPI_CLEANUP_FILES:=docs test .openapi-generator .gitignore .travis.yml git_push.sh .openapi-generator-ignore api

all: package

evalate-linux-build-privilege:
	sudo ./$(BUILD_OUTPUT_DIR)/setup-privileges.sh

install-all-go-cmds:
	go install fyne.io/fyne/v2/cmd/fyne@latest # 2.5.4
	go install github.com/zxmfke/swagger2openapi3/cmd/swag2op@latest # 0.1.1

generate-gui-resources:
	cd cmd/gui && go generate

generate-api-doc:
	mkdir -p internal/swagger-doc/openapi
	# replace to swag2op get openapi 3.0 spec for better generator experience
	# go run github.com/swaggo/swag/cmd/swag@latest init -g ./cmd/server/main.go -o internal/swagger-doc
	swag2op init -g ./cmd/server/main.go \
		-o internal/swagger-doc/ \
		-openo internal/swagger-doc/openapi

generate-api-client: generate-api-doc
	rm -rf $(API_CLIENT_OUTPUT_DIR)
	mkdir -p $(API_CLIENT_OUTPUT_DIR)
	uvx openapi-generator-cli generate -i internal/swagger-doc/openapi/swagger.yaml \
		-g go -o $(API_CLIENT_OUTPUT_DIR) \
		--skip-validate-spec \
		--additional-properties=disallowAdditionalPropertiesIfNotPresent=false,enumClassPrefix=true,structPrefix=true,generateInterfaces=true,isGoSubmodule=true,withGoMod=false,packageName=$(API_CLIENT_NAME)
	rm -rf $(addprefix $(API_CLIENT_OUTPUT_DIR)/,$(OPENAPI_CLEANUP_FILES))

generate-openapi-client:
	rm -rf $(OPENAPI_CLIENT_OUTPUT_DIR)
	mkdir -p $(OPENAPI_CLIENT_OUTPUT_DIR)
	uvx openapi-generator-cli generate -i openapi/bookxnote-local-ocr.yaml \
		-g go -o $(OPENAPI_CLIENT_OUTPUT_DIR) \
		--skip-validate-spec \
		--additional-properties=disallowAdditionalPropertiesIfNotPresent=false,enumClassPrefix=true,structPrefix=true,generateInterfaces=true,isGoSubmodule=true,withGoMod=false,packageName=$(OPENAPI_CLIENT_NAME)
	rm -rf $(addprefix $(OPENAPI_CLIENT_OUTPUT_DIR)/,$(OPENAPI_CLEANUP_FILES))

generate: generate-gui-resources generate-api-doc generate-api-client generate-openapi-client
	@echo "generated"

build: clean
	@echo "building..."
	@mkdir -p $(BUILD_OUTPUT_DIR)

	@echo "building CLI..."
	go build -o $(BUILD_OUTPUT_DIR)/$(CLI_NAME)$(BINARY_SUFFIX) ./cmd/cli/...

	@echo "building GUI..."
	go build -o $(BUILD_OUTPUT_DIR)/$(GUI_NAME)$(BINARY_SUFFIX) ./cmd/gui/...

	@echo "building server..."
	go build -o $(BUILD_OUTPUT_DIR)/$(SERVER_NAME)$(BINARY_SUFFIX) ./cmd/server/...

	@echo "building certgen..."
	go build -o $(BUILD_OUTPUT_DIR)/$(CERTGEN_NAME)$(BINARY_SUFFIX) ./cmd/certgen/...

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
		7z.exe a -tzip $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(APP_NAME)-$(PLATFORM_SUFFIX); \
	else \
		7z a -ttar $(ARCHIVE_NAME)$(ARCHIVE_EXT) $(APP_NAME)-$(PLATFORM_SUFFIX); \
	fi
	@echo "done: $(ARCHIVE_NAME)$(ARCHIVE_EXT)"

dev-build: generate prepare evalate-linux-build-privilege
	cp artifact/dev.config.yml $(BUILD_OUTPUT_DIR)/config.yml
	@echo "Done"

dev-run-gui: dev-build
	./build/bookxnote-local-ocr-linux-amd64/bookxnote-local-ocr-gui

clean:
	@echo "clean..."
	@rm -rf $(BUILD_ROOT)
