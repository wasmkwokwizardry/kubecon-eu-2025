package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
)

// FilterInput is the input to the filter scheduler plugin function.
type FilterInput struct {
	Pod      v1.Pod   `json:"pod"`
	NodeInfo NodeInfo `json:"node_info"`
}

// NodeInfo contains information about a node.
type NodeInfo struct {
	Node        v1.Node                      `json:"node"`
	ImageStates map[string]ImageStateSummary `json:"image_states,omitempty"`
}

// ImageStateSummary contains information about an image.
type ImageStateSummary struct {
	// Size of the image
	Size uint64 `json:"size"`
	// Used to track how many nodes have this image
	NumNodes uint32 `json:"num_nodes"`
}

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
		"../guest-tinygo/plugin.wasm",
		"../guest-rust/target/wasm32-unknown-unknown/release/filter_json_guest_rust.wasm",
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
			input := FilterInput{
				Pod: v1.Pod{
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
				NodeInfo: NodeInfo{
					Node: v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: nodeName,
						},
					},
				},
			}

			data, err := json.Marshal(input)
			if err != nil {
				slog.Error("Failed to marshal input", "error", err)
				continue
			}

			exit, out, err := plugin.Call("filter", data)
			if err != nil {
				slog.Error("Filter failed", "exit", exit, "error", err)
				continue
			}

			var status api.Status
			if err := json.Unmarshal(out, &status); err != nil {
				slog.Error("Failed to unmarshal status", "error", err)
				continue
			}

			slog.Info("Filter returned", "status", status)
		}
	}
}
