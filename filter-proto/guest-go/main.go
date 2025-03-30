package main

import (
	"fmt"
	"regexp"

	"github.com/extism/go-pdk"
	"google.golang.org/protobuf/proto"

	"filter-proto-host/protos/filter"
)

const (
	regexAnnotationKey = "scheduler.wasmkwokwizardry.io/regex"
)

//go:wasmexport filter
func Filter() int32 {
	input := new(filter.FilterInput)
	if err := proto.Unmarshal(pdk.Input(), input); err != nil {
		pdk.SetError(err)
		return -1
	}

	output := _filter(input)

	data, err := proto.Marshal(output)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.Output(data)

	return 0
}

func _filter(input *filter.FilterInput) *filter.Status {
	pattern, ok := input.GetPod().GetAnnotations()[regexAnnotationKey]
	if !ok {
		return &filter.Status{
			Code:   filter.StatusCode_SUCCESS.Enum(),
			Reason: proto.String("no regex annotation found"),
		}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return &filter.Status{
			Code:   filter.StatusCode_ERROR.Enum(),
			Reason: proto.String(fmt.Sprintf("failed to compile regex %q: %s", pattern, err)),
		}
	}

	nodeName := input.GetNodeInfo().GetNode().GetName()

	if regex.MatchString(nodeName) {
		return &filter.Status{
			Code:   filter.StatusCode_SUCCESS.Enum(),
			Reason: proto.String(fmt.Sprintf("node %q matches regex %q", nodeName, pattern)),
		}
	}

	return &filter.Status{
		Code:   filter.StatusCode_UNSCHEDULABLE.Enum(),
		Reason: proto.String(fmt.Sprintf("node %q does not match regex %q", nodeName, pattern)),
	}
}

func main() {}
