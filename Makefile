GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o
BIN=bin/sub2clash
SOURCE=.

docker:
	$(GOBUILD) $(BIN) $(SOURCE)