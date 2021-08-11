FROM --platform=linux/arm64 golang:alpine as builder

RUN apk add --no-cache make git

WORKDIR /sub2clash-src
COPY . /sub2clash-src
RUN export GOPROXY=https://goproxy.io,direct && \
    go mod download && \
    make docker && \
    mv ./bin/sub2clash /sub2clash

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/sub2clash"

COPY --from=builder /sub2clash /

ENTRYPOINT ["/sub2clash", "-d", "/configs"]
EXPOSE 80