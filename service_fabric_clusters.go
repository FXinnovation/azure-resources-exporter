package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/prometheus/common/log"
)

const serviceFabricClustersResourceType = "Microsoft.ServiceFabric/clusters"

// ServiceFabricClustersClient is the client implementation to ServiceFabricClusters API
type ServiceFabricClustersClient struct {
	Session   *AzureSession
	Client    *servicefabric.ClustersClient
	Resources Resources
}

// ServiceFabricClusters client interface
type ServiceFabricClusters interface {
	GetServiceFabricClusters() (*[]servicefabric.Cluster, error)
	GetSubscriptionID() string
}

// NewServiceFabricClusters returns a new ServiceFabricClusters client
func NewServiceFabricClusters(session *AzureSession) ServiceFabricClusters {

	client := servicefabric.NewClustersClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &ServiceFabricClustersClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (sc *ServiceFabricClustersClient) GetSubscriptionID() string {
	return sc.Session.SubscriptionID
}

// GetServiceFabricClusters fetch ServiceFabricClusters with state
func (sc *ServiceFabricClustersClient) GetServiceFabricClusters() (*[]servicefabric.Cluster, error) {
	var sfcList []servicefabric.Cluster

	resources, err := sc.Resources.GetResources(serviceFabricClustersResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping service fabric cluster: %s", err)
			continue
		}

		sfc, err := sc.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		sfcList = append(sfcList, sfc)
	}

	return &sfcList, nil
}
