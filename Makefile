GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

build:
	go mod tidy
	$(GOBUILD) -o bin/sub2clash .