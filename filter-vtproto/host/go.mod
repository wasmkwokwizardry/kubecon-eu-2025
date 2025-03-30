module filter-vtproto-host

go 1.24.1

require (
	github.com/extism/go-sdk v1.7.1
	github.com/planetscale/vtprotobuf v0.6.0
	github.com/tetratelabs/wazero v1.9.0
	google.golang.org/protobuf v1.36.6
	k8s.io/api v0.30.4
	k8s.io/apimachinery v0.30.4
)

require (
	github.com/dylibso/observe-sdk/go v0.0.0-20240819160327-2d926c5d788a // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20240805132620-81f5be970eca // indirect
	github.com/tetratelabs/wabin v0.0.0-20230304001439-f6f874872834 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
)

replace (
	k8s.io/api => ./protos/k8s.io/api
	k8s.io/apimachinery => ./protos/k8s.io/apimachinery
)
