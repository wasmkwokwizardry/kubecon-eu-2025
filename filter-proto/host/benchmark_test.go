package main

import (
	"fmt"
	"os"
	"testing"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"filter-proto-host/protos/filter"
)

var plugins = map[string]string{
	"Go":   "../guest-go/plugin.wasm",
	"Rust": "../guest-rust/target/wasm32-unknown-unknown/release/filter_proto_guest_rust.wasm",
}

func BenchmarkPluginFilterSequential(b *testing.B) {
	ctx := b.Context()

	cache, err := wazero.NewCompilationCacheWithDir("/tmp/wazero-cache")
	if err != nil {
		b.Fatalf("Failed to create compilation cache: %v", err)
	}
	defer cache.Close(ctx)

	config := extism.PluginConfig{
		EnableWasi:    true,
		ModuleConfig:  wazero.NewModuleConfig().WithArgs(os.Args[0]),
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(cache),
	}

	for lang, pluginPath := range plugins {
		b.Run(lang, func(b *testing.B) {
			manifest := extism.Manifest{
				Wasm: []extism.Wasm{
					extism.WasmFile{
						Path: pluginPath,
					},
				},
			}

			plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
			if err != nil {
				b.Fatalf("Failed to initialize plugin: %v", err)
			}

			// Run the benchmark
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Prepare input data without the timer
				b.StopTimer()
				input := &filter.FilterInput{
					Pod: &v1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name: fmt.Sprintf("pod-%d", i),
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
								Name: fmt.Sprintf("kubecon-%d", i),
							},
						},
					},
				}
				b.StartTimer()

				data, err := proto.Marshal(input)
				if err != nil {
					b.Fatalf("Failed to marshal input: %v", err)
				}

				exit, out, err := plugin.Call("filter", data)
				if err != nil {
					b.Fatalf("Filter failed: exit=%d, error=%v", exit, err)
				}

				status := new(filter.Status)
				if err := proto.Unmarshal(out, status); err != nil {
					b.Fatalf("Failed to unmarshal status: %v", err)
				}

				if status.GetCode() != filter.StatusCode_SUCCESS {
					b.Fatalf("Filter returned unexpected status: %v", status)
				}
			}
		})
	}
}

func BenchmarkPluginFilterParallel(b *testing.B) {
	ctx := b.Context()

	cache, err := wazero.NewCompilationCacheWithDir("/tmp/wazero-cache")
	if err != nil {
		b.Fatalf("Failed to create compilation cache: %v", err)
	}
	defer cache.Close(ctx)

	config := extism.PluginConfig{
		EnableWasi:    true,
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(cache),
	}

	for lang, pluginPath := range plugins {
		b.Run(lang, func(b *testing.B) {
			manifest := extism.Manifest{
				Wasm: []extism.Wasm{
					extism.WasmFile{
						Path: pluginPath,
					},
				},
			}

			compiledPlugin, err := extism.NewCompiledPlugin(ctx, manifest, config, []extism.HostFunction{})
			if err != nil {
				b.Fatalf("Failed to initialize plugin: %v", err)
			}

			// Run the benchmark
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				plugin, err := compiledPlugin.Instance(ctx, extism.PluginInstanceConfig{
					ModuleConfig: wazero.NewModuleConfig().WithArgs(os.Args[0]),
				})
				if err != nil {
					b.Fatalf("Failed to create plugin instance: %v", err)
				}
				defer plugin.Close(ctx)

				for pb.Next() {
					input := &filter.FilterInput{
						Pod: &v1.Pod{
							ObjectMeta: metav1.ObjectMeta{
								Name: "pod-1",
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
									Name: "kubecon-1",
								},
							},
						},
					}

					data, err := proto.Marshal(input)
					if err != nil {
						b.Fatalf("Failed to marshal input: %v", err)
					}

					exit, out, err := plugin.Call("filter", data)
					if err != nil {
						b.Fatalf("Filter failed: exit=%d, error=%v", exit, err)
					}

					status := new(filter.Status)
					if err := proto.Unmarshal(out, status); err != nil {
						b.Fatalf("Failed to unmarshal status: %v", err)
					}

					if status.GetCode() != filter.StatusCode_SUCCESS {
						b.Fatalf("Filter returned unexpected status: %v", status)
					}
				}
			})
		})
	}
}
