package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/prometheus/common/log"
)

const applicationGatewaysResourceType = "Microsoft.Network/applicationGateways"

// ApplicationGatewaysClient is the client implementation to ApplicationGateways API
type ApplicationGatewaysClient struct {
	Session   *AzureSession
	Client    *network.ApplicationGatewaysClient
	Resources Resources
}

// ApplicationGateways client interface
type ApplicationGateways interface {
	GetApplicationGateways() (*[]network.ApplicationGateway, error)
	GetSubscriptionID() string
}

// NewApplicationGateways returns a new ApplicationGateways client
func NewApplicationGateways(session *AzureSession) ApplicationGateways {

	client := network.NewApplicationGatewaysClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &ApplicationGatewaysClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (ac *ApplicationGatewaysClient) GetSubscriptionID() string {
	return ac.Session.SubscriptionID
}

// GetApplicationGateways fetch ApplicationGateways with operational state
func (ac *ApplicationGatewaysClient) GetApplicationGateways() (*[]network.ApplicationGateway, error) {
	var agList []network.ApplicationGateway

	resources, err := ac.Resources.GetResources(applicationGatewaysResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping application gateway: %s", err)
			continue
		}

		ag, err := ac.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		agList = append(agList, ag)
	}

	return &agList, nil
}
