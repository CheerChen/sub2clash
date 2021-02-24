#!/usr/bin/env bash
cd cmd/sub2clash || exit
GOOS=darwin GOARCH=arm64 go build -ldflags '-w -s'
./sub2clash -c prod