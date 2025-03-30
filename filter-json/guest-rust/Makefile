.PHONY: build build-debug test

test: build-debug
	extism call target/wasm32-unknown-unknown/debug/filter_json_guest_rust.wasm filter --input '{"pod":{"metadata":{"name":"my-pod","annotations":{"scheduler.wasmkwokwizardry.io/regex":"kubecon-.*"}}},"node_info":{"node":{"metadata":{"name":"cncf-1"}}}}' --wasi

build: k8s_api_v1.json
	cargo build --release

build-debug: k8s_api_v1.json
	cargo build

k8s_api_v1.json:
	curl -L https://github.com/yannh/kubernetes-json-schema/raw/refs/heads/master/v1.32.3-standalone-strict/_definitions.json > $@