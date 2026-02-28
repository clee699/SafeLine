# 使用官方Go镜像作为构建阶段
FROM golang:1.18 as builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建webserver
RUN CGO_ENABLED=0 GOOS=linux go build -o safeline-webserver main.go

# 使用轻量级Alpine镜像作为运行阶段
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/safeline-webserver .

# 复制配置文件
COPY config.yml .

# 复制静态资源
COPY static ./static

# 暴露端口
EXPOSE 8000

# 运行webserver
CMD ["./safeline-webserver"]