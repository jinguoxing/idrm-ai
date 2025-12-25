# IDRM 菜单管理服务 - Dockerfile
# 多阶段构建，减小镜像体积

# 阶段1: 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装必要的工具
RUN apk add --no-cache git make

# 复制go.mod和go.sum（利用Docker缓存）
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o idrm-menu-api ./api/menu.go

# 阶段2: 运行阶段
FROM alpine:latest

# 安装CA证书（如果需要调用HTTPS API）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 复制编译好的二进制文件
COPY --from=builder /build/idrm-menu-api /app/idrm-menu-api

# 复制配置文件
COPY --from=builder /build/api/etc/menu.yaml /app/menu.yaml

# 设置工作目录
WORKDIR /app

# 修改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动服务
CMD ["./idrm-menu-api", "-f", "menu.yaml"]
