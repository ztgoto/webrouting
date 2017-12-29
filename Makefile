
all: build

build: webrouting

webrouting: downlib
	go build -o ./bin/webrouting ./main

downlib:
	go get -v github.com/spf13/cobra
	go get -v github.com/valyala/fasthttp

clean:
	@rm -rf bin


# test:
# 	go test ./*
