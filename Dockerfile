FROM golang:1.18 AS build

WORKDIR /app
ADD . .

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app main.go

# 运行：使用scratch作为基础镜像
FROM scratch as prod

# 在build阶段复制时区到
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在build阶段复制可执行的go二进制文件app
COPY --from=build /app/app /
# 在build阶段复制配置文件
# COPY --from=build /app/config ./config

# 启动服务
CMD ["/app"]


