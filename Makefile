
all: build

build: webrouting

webrouting: downlib
	go build -o ./bin/webrouting ./main

downlib:
	go get -v go.uber.org/zap
	go get -v gopkg.in/yaml.v2
	go get -v github.com/spf13/cobra
	go get -v github.com/valyala/fasthttp

clean:
	@rm -rf bin


# test:
# 	go test ./*
