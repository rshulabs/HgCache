FROM golang:alpine AS builder

ENV GOPROXY https://goproxy.cn

WORKDIR /go/src/HugeCache/server

COPY go.* ./
RUN go mod tidy 
COPY . .

RUN go build -o server .

FROM alpine:latest
LABEL MAINTAINER="2505269010@qq.com"
WORKDIR /go/src/HugeCache/server
COPY --from=0 /go/src/HugeCache/server/configs/cache.yaml ./configs/cache.yaml
COPY --from=0 /go/src/HugeCache/server/run.sh ./
COPY --from=0 /go/src/HugeCache/server/server ./

EXPOSE 8791
# ENTRYPOINT ./run.sh
ENTRYPOINT ./server -port=8001 & ./server -port=8002 & ./server -port=8003 -api=1
