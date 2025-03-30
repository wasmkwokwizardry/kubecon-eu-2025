.PHONY: build build-go build-tinygo

build: build-go build-tinygo

build-go: # Big Go caching does the job for us
	GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.go.wasm main.go

build-tinygo: plugin.tinygo.wasm

plugin.tinygo.wasm: main.go go.mod go.sum
	tinygo build -target wasip1 -scheduler=none --no-debug -buildmode=c-shared -o plugin.tinygo.wasm main.go