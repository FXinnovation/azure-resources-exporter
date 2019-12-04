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
func (vc *VirtualNetworkGatewayConnectionsClient) GetSubscriptionID() string {
	return vc.Session.SubscriptionID
}

// GetVirtualNetworkGatewayConnections fetch VirtualNetworkGatewayConnections with connection status
func (vc *VirtualNetworkGatewayConnectionsClient) GetVirtualNetworkGatewayConnections() (*[]network.VirtualNetworkGatewayConnection, error) {
	var connectionList []network.VirtualNetworkGatewayConnection

	resources, err := vc.Resources.GetResources(virtualNetworkGatewayConnectionsResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceLabels(*resource.ID)
		if err != nil {
			log.Errorf("Skipping virtual network gateway connection: %s", err)
			continue
		}

		connection, err := vc.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		connectionList = append(connectionList, connection)
	}

	return &connectionList, nil
}
