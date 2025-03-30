.PHONY: build bench guest-go guest-rust test protos

test: build guest-go guest-rust
	./host

bench: guest-go guest-rust
	go test -bench=. -benchmem -benchtime=10s

build: protos
	CGO_ENABLED=0 go build -o host main.go

guest-go:
	$(MAKE) -C ../guest-go build

guest-rust:
	$(MAKE) -C ../guest-rust build

protos: protos/include protos/filter/filter.proto protos/k8s.io/api/go.mod
	protoc --go_out=protos --go-vtproto_out=protos -I . -I protos/include --go-vtproto_opt=features=marshal+unmarshal+size \
		protos/filter/filter.proto \
		k8s.io/api/core/v1/generated.proto \
		k8s.io/apimachinery/pkg/api/resource/generated.proto \
		k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto \
		k8s.io/apimachinery/pkg/runtime/generated.proto \
		k8s.io/apimachinery/pkg/runtime/schema/generated.proto \
		k8s.io/apimachinery/pkg/util/intstr/generated.proto

protos/k8s.io/api/go.mod: protos/k8s.io/apimachinery/go.mod
	mkdir -p protos/k8s.io/api
	cd protos/k8s.io/api && go mod init k8s.io/api && go mod edit -replace k8s.io/apimachinery=../apimachinery && go mod tidy

protos/k8s.io/apimachinery/go.mod:
	mkdir -p protos/k8s.io/apimachinery
	cd protos/k8s.io/apimachinery && go mod init k8s.io/apimachinery && go mod tidy

protos/include:
	go mod download k8s.io/api@v0.30.4 k8s.io/apimachinery@v0.30.4
	mkdir -p protos/include/k8s.io
	ln -s `go env GOMODCACHE`/k8s.io/api@v0.30.4 protos/include/k8s.io/api
	ln -s `go env GOMODCACHE`/k8s.io/apimachinery@v0.30.4 protos/include/k8s.io/apimachinery