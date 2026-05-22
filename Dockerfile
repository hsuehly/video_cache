# 基础镜像，基于golang的alpine镜像构建--编译阶段
FROM golang:alpine AS builder
# 全局工作目录
WORKDIR /build

COPY . .
# 下载依赖包
ENV GO111MODULE=on
#  用于代理下载go项目依赖的包
#ENV GOPROXY https://goproxy.cn,direct
RUN go mod download
# 把运行Dockerfile文件的当前目录所有文件复制到目标目录
# 环境变量
# 编译，关闭CGO，防止编译后的文件有动态链接，而alpine镜像里有些c库没有，直接没有文件的错误
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/app /build/cmd/main.go
RUN #GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /build/app
# 使用alpine这个轻量级镜像为基础镜像--运行阶段
FROM alpine AS runner
# 全局工作目录
WORKDIR /app

# 添加证书
COPY --from=builder /build/certs/ca_bundle.crt /usr/local/share/ca-certificates/

# 安装必要的软件包，设置时区，并清理缓存
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    update-ca-certificates && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
# 复制编译阶段编译出来的运行文件到目标目录
COPY --from=builder /build/app /app/app
#COPY public/404.html /app/public
# 将时区设置为东八区
#RUN echo "https://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories \
#    && echo "https://mirrors.aliyun.com/alpine/v3.8/community/" >> /etc/apk/repositories \
#    && apk add --no-cache tzdata \
#    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
#    && echo Asia/Shanghai > /etc/timezone \
#    && apk del tzdata
# 需暴露的端口
EXPOSE 8080
# 可外挂的目录
# docker run命令触发的真实命令(相当于直接运行编译后的可运行文件)
ENTRYPOINT ["./app"]