# C2V2 Tools - 开发者工具箱

一款轻量、高效的在线开发者工具集，专注于隐私保护和极致用户体验。

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 特性

| 工具 | 功能 |
|------|------|
| **JSON** | 格式化、压缩、验证，转换为 Go Struct / YAML |
| **HTML** | 美化、压缩、转义/反转义，实时客户端处理 |
| **CSS** | 美化、压缩、净化（每规则一行） |
| **Base64** | 编码、解码文本数据 |

- 🌐 **多语言**：中英文完整支持
- 🔒 **隐私优先**：所有处理在浏览器本地完成
- 🔍 **SEO 优化**：Sitemap、JSON-LD、Open Graph
- 📱 **响应式设计**：适配桌面和移动端

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 克隆仓库
git clone https://github.com/your-repo/c2v2-tools.git
cd c2v2-tools

# 配置环境变量
cp .env.example .env
vim .env  # 修改 DOMAIN 为您的域名

# 一键部署/更新
./deploy.sh
```

**其他命令：**
```bash
./deploy.sh --stop    # 停止服务
./deploy.sh --logs    # 查看日志
./deploy.sh --status  # 查看状态
```

### 本地开发

```bash
# 安装依赖
go mod tidy

# 运行
go run cmd/server/main.go

# 访问 http://localhost:5006
```

## 🛠️ 技术栈

- **后端**: Go + Gin
- **前端**: HTMX + AlpineJS + Tailwind CSS
- **部署**: Docker + Docker Compose

## 📁 项目结构

```
├── cmd/server/        # 入口
├── internal/
│   ├── app/           # 路由、配置
│   ├── middleware/    # 中间件（i18n、缓存、安全）
│   ├── pkg/           # 公共包（渲染、国际化）
│   └── tools/         # 工具实现
├── locales/           # 翻译文件
├── templates/         # HTML 模板
├── static/            # 静态资源
├── Dockerfile         # 容器构建
├── docker-compose.yml # 容器编排
└── deploy.sh          # 部署脚本
```

## 🔧 环境变量

| 变量 | 默认值 | 说明 |
|------|-------|------|
| `DOMAIN` | `http://localhost:5006` | 网站域名 |
| `PORT` | `5006` | 服务端口 |
| `SUPPORTED_LANGS` | `en,zh` | 支持的语言 |
| `DEFAULT_LANG` | `en` | 默认语言 |

## 📄 License

MIT License
