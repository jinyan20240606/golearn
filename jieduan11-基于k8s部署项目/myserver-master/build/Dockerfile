FROM golang:1.24-alpine3.21 as builder

WORKDIR /data/www/myserver
#设置变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
#复制
COPY go.mod .
COPY go.sum .
#下载go mod 相关
RUN go mod download
#子目录迭代复制
COPY . .
#查看env
RUN go env
#构建，瘦身-ldflags= "-s -w"
RUN go build -ldflags "-s -w" -o cmd/server cmd/server.go

FROM alpine:3.21

WORKDIR /data/www/myserver
#设置时区
ENV TZ=Asia/Shanghai

COPY --from=builder /data/www/myserver/cmd/server .
CMD ["./server"]
