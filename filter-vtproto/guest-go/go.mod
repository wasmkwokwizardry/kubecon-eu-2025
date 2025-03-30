module filter-vtproto-guest-go

go 1.24.1

require (
	filter-vtproto-host v0.0.0-00010101000000-000000000000
	github.com/extism/go-pdk v1.0.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/planetscale/vtprotobuf v0.6.0 // indirect
	k8s.io/api v0.30.4 // indirect
	k8s.io/apimachinery v0.30.4 // indirect
)

replace (
	filter-vtproto-host => ../host
	k8s.io/api => ../host/protos/k8s.io/api
	k8s.io/apimachinery => ../host/protos/k8s.io/apimachinery
)
