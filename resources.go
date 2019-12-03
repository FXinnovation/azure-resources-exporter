package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
)

// ResourcesClient is the client implementation to VirtualMachines API
type ResourcesClient struct {
	Session *AzureSession
	Client  *resources.Client
}

// Resources client interface
type Resources interface {
	GetResources(resourceType string) (*[]resources.GenericResource, error)
}

// NewResources returns a new Ressources client
func NewResources(session *AzureSession) Resources {
	client := resources.NewClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	return &ResourcesClient{
		Session: session,
		Client:  &client,
	}
}

// GetResources return resources by type and tags
func (rc *ResourcesClient) GetResources(resourceType string) (*[]resources.GenericResource, error) {
	filter := fmt.Sprintf("resourceType eq '%s'", resourceType)
	return rc.list(filter)
}

func (rc *ResourcesClient) list(filter string) (*[]resources.GenericResource, error) {
	var resourceList []resources.GenericResource

	for it, err := rc.Client.ListComplete(context.Background(), filter, "", nil); it.NotDone(); err = it.Next() {
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, it.Value())
	}

	return &resourceList, nil
}
