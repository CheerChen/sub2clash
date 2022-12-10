GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o

docker:
	go env -w GOPROXY="https://goproxy.cn,direct" && go env -w GOSUMDB=sum.golang.google.cn
	go mod tidy
	$(GOBUILD) bin/sub2clash .