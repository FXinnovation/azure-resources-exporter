package main

import (
	"testing"
)

func TestNewVirtualNetworkGatewayConnections_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	virtualNetworkGatewayConnections := NewVirtualNetworkGatewayConnections(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if virtualNetworkGatewayConnections.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", virtualNetworkGatewayConnections.GetSubscriptionID(), want)
	}
}
