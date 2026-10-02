# hengshui-tablet-video-rpc 后端镜像（多阶段构建）
# 构建: docker build -t myblog-backend:latest -f Dockerfile .
# 说明: 配置不烧进镜像，运行时通过挂载 /app/config.yaml 注入
FROM golang:1.22.5-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myblog-server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/myblog-server .
EXPOSE 8080
CMD ["./myblog-server"]
