package main

import (
	"testing"
)

func TestNewServiceFabricClusters_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	serviceFabricClusters := NewServiceFabricClusters(session)

	if serviceFabricClusters.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", serviceFabricClusters.GetSubscriptionID(), want)
	}
}
