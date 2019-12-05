package main

import "strings"

import "errors"

var (
	// resource component positions in a ResourceURL
	resourceGroupPosition   = 4
	resourceNamePosition    = 8
	subResourceNamePosition = 10
)

// ParseResourceLabels - Returns resource labels for a given resource ID.
func ParseResourceLabels(resourceID string) (map[string]string, error) {
	labels := make(map[string]string)
	resource := strings.Split(resourceID, "/")

	if len(resource) < resourceNamePosition+1 {
		return nil, errors.New("Error parsing resource ID, expected pattern is not matched for " + resourceID)
	}

	labels["resource_group"] = resource[resourceGroupPosition]
	labels["resource_name"] = resource[resourceNamePosition]
	if len(resource) > subResourceNamePosition {
		labels["sub_resource_name"] = resource[subResourceNamePosition]
	}
	return labels, nil
}
