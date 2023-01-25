FROM golang:1.18 as builder
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on
ENV GOCACHE=/go/pkg/.cache/go-build

ADD . /skyey2-flood
WORKDIR /skyey2-flood
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /skyey2-flood/skyey2-flood /skyey2-flood/main

FROM alpine:3.6 as alpine
RUN apk update && \
    apk add -U --no-cache ca-certificates tzdata

FROM alpine:3.6
MAINTAINER zhuweitung
LABEL maintainer="zhuweitung" \
    email="zhuweitung@foxmail.com"

ENV TZ="Asia/Shanghai"

COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /skyey2-flood/skyey2-flood /skyey2-flood/skyey2-flood
COPY --from=builder /skyey2-flood/config /skyey2-flood/config

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /skyey2-flood
CMD ["./skyey2-flood"]
