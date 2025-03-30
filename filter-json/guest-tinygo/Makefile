.PHONY: build test

test: build
	extism call plugin.wasm filter --input '{"pod":{"metadata":{"name":"my-pod","annotations":{"scheduler.wasmkwokwizardry.io/regex":"kubecon-.*"}}},"node_info":{"node":{"metadata":{"name":"cncf-1"}}}}' --wasi

build: plugin.wasm

plugin.wasm: main.go go.mod go.sum
	tinygo build -target wasip1 -scheduler=none --no-debug -buildmode=c-shared -o plugin.wasm main.go