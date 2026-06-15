# 第一个阶段：专门用来做编译
FROM golang:1.19.4-alpine3.17 AS builder

WORKDIR /build
# 关闭 CGO，静态编译（不依赖系统 glibc，alpine 可直接运行）
ENV CGO_ENABLED 0
ENV GOOS linux
# 配置国内 Go 模块代理，加速依赖下载
ENV GOPROXY https://goproxy.cn,direct
# 先拷贝依赖文件，利用 Docker 缓存：依赖不变就不重复 download
COPY go.mod .
COPY go.sum .
RUN go mod download
# 拷贝项目全部源码到容器 /build
COPY . .
# 打印 Go 环境信息（调试用，正式镜像可删掉）
RUN go env
# 编译：-ldflags="-s -w" 剥离符号表、调试信息，精简二进制体积
# 输出可执行文件名为 hello，入口文件 main.go
RUN go build  -ldflags="-s -w" -o hello main.go

# 多阶段构建，写多个from，可以把上个阶段的容器镜像结果放到这个容器中使用
# 第二个阶段：专门用来输出，依赖第一阶段结果

FROM alpine

WORKDIR /build
# 设置时区
ENV TZ Asia/Shanghai

COPY --from=builder /build/hello .

CMD ["./hello"]