GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

build-amd64:
	go mod tidy
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/sub2clash .
build-arm64:
	go mod tidy
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/sub2clash .