FROM golang:1.13-alpine AS go-builder
ENV GO111MODULE=on
WORKDIR /go

RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.aliyun.com/g' -i~ /etc/apk/repositories

# add git for go modules
RUN apk update && apk add make git upx
RUN go env -w GOPROXY=https://goproxy.io

COPY . .

RUN make
RUN upx ./bin/broker

FROM alpine

LABEL MAINTAINER="kaisawind <wind.kaisa@gmail.com>"

WORKDIR /
COPY --from=go-builder /go/bin/ /usr/local/bin/

CMD ["broker --grpc-host 0.0.0.0 --grpc-port 6653"]
