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

protos: protos/include protos/filter/filter.proto
	protoc --go_out=protos -I . -I protos/include \
		protos/filter/filter.proto

protos/include:
	go mod download k8s.io/api@v0.30.4 k8s.io/apimachinery@v0.30.4
	mkdir -p protos/include/k8s.io
	ln -s `go env GOMODCACHE`/k8s.io/api@v0.30.4 protos/include/k8s.io/api
	ln -s `go env GOMODCACHE`/k8s.io/apimachinery@v0.30.4 protos/include/k8s.io/apimachinery