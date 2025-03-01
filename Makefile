# 项目名称和主要变量
PROJECT_NAME := logsnap
MAIN_PACKAGE := ./main.go
OUTPUT_NAME := logsnap
BINARY_NAME := logsnap
GO := go

# 版本信息
VERSION := $(shell cat version.json | jq -r '.version')
BUILD_TIME := $(shell date +%FT%T%z)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
UPLOAD_CONFIG_URL ?= "http://127.0.0.1:5500/config.json"
DOWNLOAD_CONFIG_URL ?= "http://127.0.0.1:5500/version.json"

# Go 构建标志
LDFLAGS := -ldflags "-X logsnap/version.Version=$(VERSION) -X logsnap/version.BuildTime=$(BUILD_TIME) -X logsnap/version.GitCommit=$(GIT_COMMIT) -X logsnap/constants.UploadConfigURL=$(UPLOAD_CONFIG_URL) -X logsnap/constants.DownloadConfigURL=$(DOWNLOAD_CONFIG_URL)"
BUILD_FLAGS := -v

# 优化构建标志 (用于减小二进制文件大小)
OPTIMIZE_LDFLAGS := -ldflags "-s -w -X logsnap/version.Version=$(VERSION) -X logsnap/version.BuildTime=$(BUILD_TIME) -X logsnap/version.GitCommit=$(GIT_COMMIT) -X logsnap/constants.UploadConfigURL=$(UPLOAD_CONFIG_URL) -X logsnap/constants.DownloadConfigURL=$(DOWNLOAD_CONFIG_URL)"
OPTIMIZE_BUILD_FLAGS := -trimpath

# 目标目录
BIN_DIR := ./bin
DIST_DIR := ./dist

# 测试覆盖率阈值（百分比）
COVERAGE_THRESHOLD := 75

# 交叉编译相关变量
MAIN_PATH := ./                  # 修改为项目根目录

# 默认目标
.PHONY: all
all: clean build

# 构建应用
.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(OUTPUT_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BIN_DIR)/$(OUTPUT_NAME)"

# 运行应用
.PHONY: run
run: build
	@echo "Running $(PROJECT_NAME)..."
	@$(BIN_DIR)/$(OUTPUT_NAME)

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(DIST_DIR)
	$(GO) test -coverprofile=$(DIST_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(DIST_DIR)/coverage.out -o $(DIST_DIR)/coverage.html
	@echo "Coverage report generated at $(DIST_DIR)/coverage.html"
	@echo "Checking coverage threshold ($(COVERAGE_THRESHOLD)%)..."
	@coverage=$$($(GO) tool cover -func=$(DIST_DIR)/coverage.out | grep total: | grep -Eo '[0-9]+\.[0-9]+' | awk '{ printf "%.2f", $$1 }'); \
	echo "Total coverage: $$coverage%"; \
	if echo "$$coverage < $(COVERAGE_THRESHOLD)" | bc -l | grep -q 1; then \
		echo "❌ Coverage is below threshold: $$coverage% < $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "✅ Coverage meets threshold: $$coverage% >= $(COVERAGE_THRESHOLD)%"; \
	fi

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)
	@echo "Clean complete"

# 创建发布包
.PHONY: dist
dist: clean build-small
	@echo "Creating distribution package..."
	@mkdir -p $(DIST_DIR)
	@cp $(BIN_DIR)/$(OUTPUT_NAME) $(DIST_DIR)/
	@cp README_zh.md $(DIST_DIR)/README.md
	@cd $(DIST_DIR) && tar -czf $(PROJECT_NAME)-$(VERSION).tar.gz $(OUTPUT_NAME) README.md
	@echo "Distribution package created: $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION).tar.gz"

# 打包目标平台的二进制文件
.PHONY: package
package: build-all
	@echo "Packaging binaries..."
	@mkdir -p $(DIST_DIR)
	@cd $(BIN_DIR)/linux-amd64 && zip -9 ../../$(DIST_DIR)/logsnap-linux-amd64.zip logsnap
	@cd $(BIN_DIR)/windows-amd64 && zip -9 ../../$(DIST_DIR)/logsnap-windows-amd64.zip logsnap.exe
	@echo "Package complete:"
	@echo " - $(DIST_DIR)/logsnap-linux-amd64.zip"
	@echo " - $(DIST_DIR)/logsnap-windows-amd64.zip"

# 帮助信息
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build          - 构建应用"
	@echo "  make build-small    - 构建优化版本 (更小的二进制文件)"
	@echo "  make build-tiny     - 构建并使用UPX压缩的极小版本"
	@echo "  make run            - 构建并运行应用"
	@echo "  make deps           - 安装依赖"
	@echo "  make test           - 运行测试"
	@echo "  make test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  make lint           - 运行代码检查"
	@echo "  make fmt            - 格式化代码"
	@echo "  make clean          - 清理构建产物"
	@echo "  make dist           - 创建发布包"
	@echo "  make package        - 打包各平台的二进制文件为 zip 格式"
	@echo "  make help           - 显示帮助信息"

# 编译 Linux x86_64 版本
.PHONY: linux-amd64
linux-amd64:
	@echo "Building Linux amd64 version..."
	@mkdir -p $(BIN_DIR)/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(OPTIMIZE_BUILD_FLAGS) $(OPTIMIZE_LDFLAGS) \
		-o "$(BIN_DIR)/linux-amd64/$(BINARY_NAME)" $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/linux-amd64/$(BINARY_NAME)"

# 编译 Windows x86_64 版本
.PHONY: windows-amd64
windows-amd64:
	@echo "Building Windows amd64 version..."
	@mkdir -p $(BIN_DIR)/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(OPTIMIZE_BUILD_FLAGS) $(OPTIMIZE_LDFLAGS) \
		-o "$(BIN_DIR)/windows-amd64/$(BINARY_NAME).exe" $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/windows-amd64/$(BINARY_NAME).exe"

# 编译所有平台
.PHONY: build-all
build-all: clean linux-amd64 windows-amd64
	@echo "All platform builds completed" 