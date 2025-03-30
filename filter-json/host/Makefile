.PHONY: build bench guest-go guest-tinygo guest-rust test

test: build
	./host

bench: guest-go guest-tinygo guest-rust
	go test -bench=. -benchmem -benchtime=10s

build: host guest-go guest-tinygo guest-rust

guest-go:
	$(MAKE) -C ../guest-go build

guest-tinygo:
	$(MAKE) -C ../guest-tinygo build

guest-rust:
	$(MAKE) -C ../guest-rust build

host: main.go
	CGO_ENABLED=0 go build -o host main.go