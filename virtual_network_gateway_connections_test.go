package main

import (
	"testing"
)

func TestNewVirtualNetworkGatewayConnections_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	session, err := NewAzureSession(wantSubscriptionID)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	virtualNetworkGatewayConnections := NewVirtualNetworkGatewayConnections(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if virtualNetworkGatewayConnections.GetSubscriptionID() != wantSubscriptionID {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", virtualNetworkGatewayConnections.GetSubscriptionID(), wantSubscriptionID)
	}
}
