.PHONY: build clean test run lint swagger migrate

# 变量
APP_NAME=ssop
BIN_DIR=./bin
API_BINARY=$(BIN_DIR)/api
MAIN_API_FILE=./cmd/api/main.go

# 环境设置
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 默认目标
all: clean build

# 构建应用
build:
	@echo "Building API..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(API_BINARY) $(MAIN_API_FILE)

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	@rm -rf $(BIN_DIR)

# 运行应用
run:
	@go run $(MAIN_API_FILE)

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 代码格式化与静态检查
lint:
	@echo "Running linters..."
	@golangci-lint run ./...
	@go fmt ./...

# 生成Swagger文档
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g $(MAIN_API_FILE) -o ./docs/swagger

# 运行数据库迁移
migrate:
	@echo "Running database migrations..."
	@go run ./scripts/migrate.go

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 帮助信息
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run code linters"
	@echo "  make swagger  - Generate Swagger documentation"
	@echo "  make migrate  - Run database migrations"
	@echo "  make deps     - Install development dependencies" 