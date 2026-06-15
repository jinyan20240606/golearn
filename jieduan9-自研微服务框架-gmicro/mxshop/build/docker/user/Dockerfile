FROM golang:1.19.4-alpine3.17 AS builder

WORKDIR /build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build  -ldflags="-s -w" -o user cmd/user/user.go

FROM alpine

WORKDIR /build

ENV TZ Asia/Shanghai

# 创建配置目录
RUN mkdir -p configs/user
# 创建日志目录，-p就是递归创建父目录
RUN mkdir -p logs

COPY --from=builder /build/user .
# 声明容器内部应用监听的端口为 8021
EXPOSE 8021
# 固定容器主程序：数组格式（Exec 模式），Docker 标准推荐写法
# 定义容器启动后第一个执行的主程序。
# 规则：
# 数组形式不会被 Shell 解析，直接运行二进制 / 可执行文件；
# ENTRYPOINT 是主进程，容器生命周期由它决定；
# 它会拼接后面 CMD 的内容作为启动参数。
ENTRYPOINT ["./user"]
# 当 Dockerfile 同时存在 ENTRYPOINT + CMD 时：
# CMD 的内容全部当作参数，追加给 ENTRYPOINT 程序
CMD ["-c", "configs/user/srv.yaml"]



# # 固定前缀：./user --env prod
# ENTRYPOINT ["./user", "--env", "prod"]
# # CMD 作为可变后缀参数
# CMD ["-c", "configs/user/srv.yaml"]
# 结果：./user --env prod -c configs/user/srv.yaml