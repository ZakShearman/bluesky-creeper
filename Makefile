# Define the directories for protobuf files and generated code
PROTO_DIR := pkg/proto
GENERATED_DIR := pkg/protobuf/generated

ifeq ($(OS),Windows_NT)
	MKDIR = mkdir $(subst /,\,$(GENERATED_DIR))
    FOR_LOOP = for %%f in ($(PROTO_DIR)\*.proto) do $(PROTOC) --go_out=$(GENERATED_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GENERATED_DIR) --go-grpc_opt=paths=source_relative -I $(PROTO_DIR) %%f

else
	MKDIR = mkdir -p $(GENERATED_DIR)
	FOR_LOOP = for proto_file in $(PROTO_DIR)/*.proto; do \
    		$(PROTOC) --go_out=$(GENERATED_DIR) --go_opt=paths=source_relative \
			--go-grpc_out=$(GENERATED_DIR) --go-grpc_opt=paths=source_relative \
			-I $(PROTO_DIR) $$proto_file; \
		done
endif

# Protobuf and Go options
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# The Go import path for generated code
GO_IMPORT_PATH := your_project_path

# Define the URL for downloading protoc (adjust version as needed)
PROTOC_VERSION := 28.3# Replace with the desired version
PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(OS)-x86_64.zip

# Default target: generate all protobuf files
all: generate

# Generate all protobuf files
generate: $(PROTO_DIR)/*.proto
	$(MKDIR)
	$(FOR_LOOP)

# Clean the generated files
clean:
	rm -rf $(GENERATED_DIR)/*.pb.go

# Install the necessary tools (protoc, protoc-gen-go, protoc-gen-go-grpc)
install_tools:
	@echo "Installing protoc..."
	@if "$(OS)"=="Windows_NT" (
		REM Install protoc for Windows
		curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v28.3/protoc-28.3-Windows_NT-x86_64.zip
		REM Unzip the protoc binary using PowerShell
		powershell -Command "Expand-Archive -Path protoc-28.3-Windows_NT-x86_64.zip -DestinationPath protoc"
		move protoc\bin\protoc.exe $(GOPATH)\bin\
		del /q protoc-28.3-Windows_NT-x86_64.zip
		del /q /f /s /q protoc
	) else (
		@if [ "$(shell echo $$OSTYPE)" = "darwin"* ]; then \
			# Install protoc for macOS \
			curl -LO $(PROTOC_URL); \
			tar -xvzf protoc-$(PROTOC_VERSION)-osx-x86_64.tar.gz -C protoc; \
			mv protoc/bin/protoc $(GOPATH)/bin/; \
			rm -rf protoc; \
		elif [ "$(shell echo $$OSTYPE)" = "linux"* ]; then \
			# Install protoc for Linux \
			curl -LO $(PROTOC_URL); \
			tar -xvzf protoc-$(PROTOC_VERSION)-linux-x86_64.tar.gz -C protoc; \
			mv protoc/bin/protoc $(GOPATH)/bin/; \
			rm -rf protoc; \
		else \
			echo "Unsupported OS"; \
			exit 1; \
		fi
	)
	# Install protoc-gen-go and protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


.PHONY: all generate clean install_tools