package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"filter-proto-host/protos/filter"
)

func main() {
	ctx := context.Background()

	// This is optional
	cache, err := wazero.NewCompilationCacheWithDir("/tmp/wazero-cache")
	if err != nil {
		log.Fatalf("Failed to create compilation cache: %v", err)
	}
	defer cache.Close(ctx)

	config := extism.PluginConfig{
		EnableWasi:    true,
		ModuleConfig:  wazero.NewModuleConfig().WithArgs(os.Args[0]), // needed for klog to init properly with a Big Go guest
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(cache),
	}

	for _, pluginPath := range []string{
		"../guest-go/plugin.wasm",
		"../guest-rust/target/wasm32-unknown-unknown/release/filter_proto_guest_rust.wasm",
	} {
		manifest := extism.Manifest{
			Wasm: []extism.Wasm{
				extism.WasmFile{
					Path: pluginPath,
				},
			},
		}

		plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})

		if err != nil {
			log.Fatalf("Failed to initialize plugin: %v", err)
		}

		for _, nodeName := range []string{"kubecon-1", "cncf-1"} {
			input := &filter.FilterInput{
				Pod: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "my-pod",
						Annotations: map[string]string{
							"scheduler.wasmkwokwizardry.io/regex": "kubecon-.*",
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{},
					},
				},
				NodeInfo: &filter.NodeInfo{
					Node: &v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: nodeName,
						},
					},
				},
			}

			data, err := proto.Marshal(input)
			if err != nil {
				slog.Error("Failed to marshal input", "error", err)
				continue
			}

			exit, out, err := plugin.Call("filter", data)
			if err != nil {
				slog.Error("Filter failed", "exit", exit, "error", err)
				continue
			}

			status := new(filter.Status)
			if err := proto.Unmarshal(out, status); err != nil {
				slog.Error("Failed to unmarshal output", "error", err)
				continue
			}

			slog.Info("Filter returned", "status", status)
		}
	}
}
