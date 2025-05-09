#!/bin/bash
ENV_DIR="/home/sinker/WorkSpace/golang/ssop/ssop_rear"
if [ -f "$ENV_DIR/.env" ]; then
  set -o allexport
  source <(grep -v '^#' "$ENV_DIR/.env" | sed '/^\s*$/d')
  set +o allexport
fi


# 获取数据库连接信息
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-"skd1701"}
DB_HOST=${DB_HOST:-10.8.118.114}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-ssop}
echo "${DB_USER}"
# 构建数据库连接字符串
CONNECTION_STRING="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"
echo "${CONNECTION_STRING}"
# 确认是否安装了 goose
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed. Please install it with 'go install github.com/pressly/goose/v3/cmd/goose@latest'"
    exit 1
fi

# 执行迁移
ACTION=${1:-"up"}
case $ACTION in
    "up")
        echo "Running migrations up..."
        goose -dir ./migrations mysql "${CONNECTION_STRING}" up
        ;;
    "down")
        echo "Rolling back last migration..."
        goose -dir ./migrations mysql "${CONNECTION_STRING}" down
        ;;
    "reset")
        echo "Resetting database (down and then up)..."
        goose -dir ./migrations mysql "${CONNECTION_STRING}" reset
        goose -dir ./migrations mysql "${CONNECTION_STRING}" up
        ;;
    "status")
        echo "Migration status..."
        goose -dir ./migrations mysql "${CONNECTION_STRING}" status
        ;;
    *)
        echo "Invalid action. Use: up, down, reset or status"
        exit 1
        ;;
esac

echo "Migration action '${ACTION}' completed."