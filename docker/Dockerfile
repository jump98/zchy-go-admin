# 第一阶段：编译
FROM golang:1.23-alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/app
WORKDIR /app

# 将代码复制到容器中
ADD go.mod .
ADD go.sum .
# 下载依赖信息
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 /go-admin
RUN go build -o /go-admin ./main.go

# 第二阶段：运行镜像
FROM alpine:3.20

ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 只复制二进制文件
COPY --from=builder /go-admin /go-admin
# 复制 config 文件夹
COPY --from=builder /app/config /config

# 执行命令
CMD ["/go-admin","server","-c", "/config/settings.prod.yml","-t","true"]



