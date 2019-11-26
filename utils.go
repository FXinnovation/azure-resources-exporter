package main

import "strings"

var (
	// resource component positions in a ResourceURL
	resourceGroupPosition   = 4
	resourceNamePosition    = 8
	subResourceNamePosition = 10
)

// ParseResourceLabels - Returns resource labels for a given resource ID.
func ParseResourceLabels(resourceID string) map[string]string {
	labels := make(map[string]string)
	resource := strings.Split(resourceID, "/")

	labels["resource_group"] = resource[resourceGroupPosition]
	labels["resource_name"] = resource[resourceNamePosition]
	if len(resource) > 13 {
		labels["sub_resource_name"] = resource[subResourceNamePosition]
	}
	return labels
}
