package main

import (
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
)

func BenchmarkNativeFilterSequential(b *testing.B) {
	b.Run("native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Prepare input data without the timer
			b.StopTimer()
			input := FilterInput{
				Pod: v1.Pod{
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
				NodeInfo: NodeInfo{
					Node: v1.Node{
						ObjectMeta: metav1.ObjectMeta{
							Name: fmt.Sprintf("kubecon-%d", i),
						},
					},
				},
			}
			b.StartTimer()

			status := Filter(&input)

			if status.Code != api.StatusCodeSuccess {
				b.Fatalf("Filter returned unexpected status: %v", status)
			}
		}
	})
}

func BenchmarkNativeFilterParallel(b *testing.B) {
	b.Run("native", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				input := FilterInput{
					Pod: v1.Pod{
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
					NodeInfo: NodeInfo{
						Node: v1.Node{
							ObjectMeta: metav1.ObjectMeta{
								Name: "kubecon-1",
							},
						},
					},
				}

				status := Filter(&input)

				if status.Code != api.StatusCodeSuccess {
					b.Fatalf("Filter returned unexpected status: %v", status)
				}
			}
		})
	})
}
