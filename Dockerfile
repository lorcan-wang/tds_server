#build 阶段
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# 拷贝所有代码
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go 

# 最小化运行镜像
FROM --platform=linux/amd64 alpine:3.19

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]