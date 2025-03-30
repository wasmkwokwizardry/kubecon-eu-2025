package main

import (
	"fmt"
	"regexp"

	"github.com/extism/go-pdk"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
)

type ImageStateSummary struct {
	// Size of the image
	Size uint64 `json:"size"`
	// Used to track how many nodes have this image
	NumNodes uint32 `json:"num_nodes"`
}

type NodeInfo struct {
	Node        v1.Node                      `json:"node"`
	ImageStates map[string]ImageStateSummary `json:"image_states"`
}

type FilterInput struct {
	Pod      v1.Pod   `json:"pod"`
	NodeInfo NodeInfo `json:"node_info"`
}

const (
	regexAnnotationKey = "scheduler.wasmkwokwizardry.io/regex"
)

//go:wasmexport filter
func Filter() int32 {
	var input FilterInput
	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return -1
	}

	output := filter(&input)

	if err := pdk.OutputJSON(output); err != nil {
		pdk.SetError(err)
		return -1
	}

	return 0
}

func filter(input *FilterInput) api.Status {
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

func main() {}
