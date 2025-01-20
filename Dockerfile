FROM golang:alpine as builder

WORKDIR /src
COPY . /src

RUN go mod tidy && \
    CGO_ENABLED=0 go build -trimpath -ldflags '-w -s' -o bin/init . && \
    mv ./bin/init /init

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/sub2clash"

COPY --from=builder /init /

ENTRYPOINT ["/init"]