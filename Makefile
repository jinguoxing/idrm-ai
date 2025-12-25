# IDRM 菜单管理模块 - Makefile
# 提供常用的构建、测试、运行命令

.PHONY: all build run test clean fmt lint deps help

# 变量定义
APP_NAME=idrm-menu-api
CMD_DIR=./api
BUILD_DIR=./bin
CONFIG_FILE=$(CMD_DIR)/etc/menu.yaml

# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOFMT=gofmt
GOIMPORTS=goimports

# 默认目标
all: fmt build

## help: 显示帮助信息
help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║          IDRM 菜单管理模块 - Make 命令帮助                      ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用命令:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""

## deps: 下载依赖
deps:
	@echo "📦 下载Go模块依赖..."
	$(GOCMD) mod download
	$(GOCMD) mod tidy

## fmt: 格式化代码
fmt:
	@echo "✨ 格式化Go代码..."
	$(GOFMT) -s -w .
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		$(GOIMPORTS) -w .; \
	fi

## lint: 代码检查
lint:
	@echo "🔍 运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint未安装，跳过"; \
	fi

## vet: 静态分析
vet:
	@echo "🔬 运行go vet静态分析..."
	$(GOCMD) vet ./...

## build: 编译项目
build: clean
	@echo "🔨 编译项目..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/menu.go
	@echo "✅ 编译完成: $(BUILD_DIR)/$(APP_NAME)"

## run: 运行服务
run:
	@echo "🚀 启动服务..."
	$(GORUN) $(CMD_DIR)/menu.go -f $(CONFIG_FILE)

## test: 运行测试
test:
	@echo "🧪 运行单元测试..."
	$(GOTEST) -v ./model/system/menu/... -cover

## test-all: 运行所有测试
test-all:
	@echo "🧪 运行所有测试..."
	$(GOTEST) -v ./... -cover

## test-cover: 生成测试覆盖率报告
test-cover:
	@echo "📊 生成测试覆盖率报告..."
	$(GOTEST) ./model/system/menu/... -coverprofile=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告已生成: coverage.html"

## benchmark: 性能测试
benchmark:
	@echo "⚡ 运行性能测试..."
	$(GOTEST) -bench=. -benchmem ./model/system/menu/...

## clean: 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

## migrate: 执行数据库迁移
migrate:
	@echo "💾 执行数据库迁移..."
	@if [ -f "deploy/migrations/20241224_create_sys_menu.sql" ]; then \
		mysql -u root -p idrm_db < deploy/migrations/20241224_create_sys_menu.sql; \
		echo "✅ 数据库迁移完成"; \
	else \
		echo "❌ 迁移文件不存在"; \
	fi

## api-gen: 生成API代码 (需要goctl)
api-gen:
	@echo "🔧 生成API代码..."
	@if command -v goctl >/dev/null 2>&1; then \
		goctl api go -api api/sys/menu.api -dir .; \
		echo "✅ API代码生成完成"; \
	else \
		echo "❌ goctl未安装，请运行: go install github.com/zeromicro/go-zero/tools/goctl@latest"; \
	fi

## install-goctl: 安装goctl工具
install-goctl:
	@echo "📦 安装goctl工具..."
	$(GOCMD) install github.com/zeromicro/go-zero/tools/goctl@latest

## install-tools: 安装开发工具
install-tools:
	@echo "🛠️  安装开发工具..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ 开发工具安装完成"

## docker-build: 构建Docker镜像
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t idrm-menu-api:latest .

## docker-run: 运行Docker容器
docker-run:
	@echo "🐳 运行Docker容器..."
	docker run -d -p 8080:8080 --name idrm-menu-api idrm-menu-api:latest

## docker-stop: 停止Docker容器
docker-stop:
	@echo "🐳 停止Docker容器..."
	docker stop idrm-menu-api || true
	docker rm idrm-menu-api || true

## check: 完整检查（fmt + vet + lint + test）
check: fmt vet lint test
	@echo "✅ 所有检查通过"

## dev: 开发模式（自动重启）
dev:
	@echo "🔥 开发模式启动..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air未安装，请运行: go install github.com/cosmtrek/air@latest"; \
		$(GORUN) $(CMD_DIR)/menu.go -f $(CONFIG_FILE); \
	fi
