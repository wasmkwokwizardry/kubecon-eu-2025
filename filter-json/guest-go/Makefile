.PHONY: build test

test: build
	extism call plugin.wasm filter --input '{"pod":{"metadata":{"name":"my-pod","annotations":{"scheduler.wasmkwokwizardry.io/regex":"kubecon-.*"}}},"node_info":{"node":{"metadata":{"name":"cncf-1"}}}}' --wasi

build: plugin.wasm

plugin.wasm: main.go
	GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.wasm main.go