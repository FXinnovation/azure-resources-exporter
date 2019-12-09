package main

import "strings"

import "errors"

var (
	// resource component positions in a ResourceURL
	resourceGroupPosition   = 4
	resourceNamePosition    = 8
	subResourceNamePosition = 10
)

// ParseResourceID - Returns resource info from a given resource ID.
func ParseResourceID(resourceID string) (map[string]string, error) {
	info := make(map[string]string)
	resource := strings.Split(resourceID, "/")

	if len(resource) < resourceNamePosition+1 {
		return nil, errors.New("Error parsing resource ID, expected pattern is not matched for " + resourceID)
	}

	info["resource_group"] = resource[resourceGroupPosition]
	info["resource_name"] = resource[resourceNamePosition]
	if len(resource) > subResourceNamePosition {
		info["sub_resource_name"] = resource[subResourceNamePosition]
	}
	return info, nil
}
