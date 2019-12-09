package main

import (
	"testing"
)

func TestNewApplicationGateways_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	applicationGateways := NewApplicationGateways(session)

	if applicationGateways.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", applicationGateways.GetSubscriptionID(), want)
	}
}
