package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/prometheus/common/log"
)

const virtualNetworkGatewayConnectionsResourceType = "Microsoft.Network/connections"

// VirtualNetworkGatewayConnectionsClient is the client implementation to VirtualNetworkGatewayConnections API
type VirtualNetworkGatewayConnectionsClient struct {
	Session   *AzureSession
	Client    *network.VirtualNetworkGatewayConnectionsClient
	Resources Resources
}

// VirtualNetworkGatewayConnections client interface
type VirtualNetworkGatewayConnections interface {
	GetVirtualNetworkGatewayConnections() (*[]network.VirtualNetworkGatewayConnection, error)
	GetSubscriptionID() string
}

// NewVirtualNetworkGatewayConnections returns a new VirtualNetworkGatewayConnections client
func NewVirtualNetworkGatewayConnections(session *AzureSession) VirtualNetworkGatewayConnections {

	client := network.NewVirtualNetworkGatewayConnectionsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &VirtualNetworkGatewayConnectionsClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (client *VirtualNetworkGatewayConnectionsClient) GetSubscriptionID() string {
	return client.Session.SubscriptionID
}

// GetVirtualNetworkGatewayConnections fetch VirtualNetworkGatewayConnections with connection status
func (client *VirtualNetworkGatewayConnectionsClient) GetVirtualNetworkGatewayConnections() (*[]network.VirtualNetworkGatewayConnection, error) {
	var connectionList []network.VirtualNetworkGatewayConnection

	ressources, err := client.Resources.GetResources(virtualNetworkGatewayConnectionsResourceType)
	if err != nil {
		return nil, err
	}

	for _, ressource := range *ressources {
		labels, err := ParseResourceLabels(*ressource.ID)
		if err != nil {
			log.Errorf("Skipping virtual network gateway connection: %s", err)
			continue
		}

		connection, err := client.Client.Get(context.Background(), labels["resource_group"], *ressource.Name)
		if err != nil {
			return nil, err
		}

		connectionList = append(connectionList, connection)
	}

	return &connectionList, nil
}
