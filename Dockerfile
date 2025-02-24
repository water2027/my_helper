# 使用官方 Go 语言镜像作为基础镜像
FROM golang:1.18-alpine

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目文件
COPY . .

# 构建 Go 应用程序
RUN go build -o main .

COPY config.yaml .

# 暴露应用程序端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]