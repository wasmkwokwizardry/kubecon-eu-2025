package main

import (
	"fmt"
	"log/slog"
	"regexp"

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

		status := Filter(&input)

		slog.Info("Filter returned", "status", status)
	}
}

const regexAnnotationKey = "scheduler.wasmkwokwizardry.io/regex"

func Filter(input *FilterInput) api.Status {
	pattern, ok := input.Pod.GetAnnotations()[regexAnnotationKey]
	if !ok {
		return api.Status{Code: api.StatusCodeSuccess, Reason: "no regex annotation found"}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return api.Status{Code: api.StatusCodeError, Reason: fmt.Sprintf("failed to compile regex %q: %s", pattern, err)}
	}

	nodeName := input.NodeInfo.Node.GetName()

	if regex.MatchString(nodeName) {
		return api.Status{Code: api.StatusCodeSuccess, Reason: fmt.Sprintf("node %q matches regex %q", nodeName, pattern)}
	}

	return api.Status{Code: api.StatusCodeUnschedulable, Reason: fmt.Sprintf("node %q does not match regex %q", nodeName, pattern)}
}
