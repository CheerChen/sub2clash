FROM golang:alpine as builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache make git

WORKDIR /sub2clash-src
COPY . /sub2clash-src

RUN make docker && \
    mv ./bin/sub2clash /sub2clash

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/sub2clash"

COPY --from=builder /sub2clash /

ENTRYPOINT ["/sub2clash"]
EXPOSE 80