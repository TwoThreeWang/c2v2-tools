#!/bin/bash
# C2V2 Tools 一键部署/更新脚本
# 用法: ./deploy.sh [--stop|--logs|--status]

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

cd "$(dirname "$0")"

# 检查依赖
command -v docker &>/dev/null || { log_error "Docker 未安装"; exit 1; }
docker compose version &>/dev/null || { log_error "Docker Compose 未安装"; exit 1; }

# 检查 .env 文件
if [ ! -f .env ]; then
    log_warn ".env 文件不存在，从示例文件创建..."
    cp .env.example .env
    log_warn "请编辑 .env 文件配置 DOMAIN 后重新运行"
    exit 0
fi

# 部署/更新
deploy() {
    # Git 更新检测
    if [ -d .git ]; then
        log_info "检查代码更新..."
        git fetch origin 2>/dev/null || true
        LOCAL=$(git rev-parse HEAD 2>/dev/null)
        REMOTE=$(git rev-parse @{u} 2>/dev/null || echo "$LOCAL")
        if [ "$LOCAL" != "$REMOTE" ]; then
            log_info "拉取最新代码..."
            git pull
        fi
    fi

    log_info "构建并启动服务..."
    docker compose up -d --build

    sleep 3
    if docker compose ps | grep -q "Up"; then
        PORT=$(grep -E "^PORT=" .env | cut -d '=' -f2 || echo "5006")
        log_info "✅ 服务运行中: http://localhost:${PORT}"
    else
        log_error "启动失败"; docker compose logs; exit 1
    fi
}

case "$1" in
    --stop)   docker compose down && log_info "服务已停止" ;;
    --logs)   docker compose logs -f ;;
    --status) docker compose ps ;;
    *)        deploy ;;
esac
