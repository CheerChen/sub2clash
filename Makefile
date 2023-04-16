GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o

docker:
	go mod tidy
	$(GOBUILD) bin/sub2clash .