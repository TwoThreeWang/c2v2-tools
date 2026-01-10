# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制模板和静态文件
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/locales ./locales
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 5006

# 生产环境模式
ENV GIN_MODE=release

# 运行应用
CMD ["./main"]
