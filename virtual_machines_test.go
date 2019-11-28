package main

import (
	"testing"
)

func TestNewVirtualMachines_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	session, err := NewAzureSession(wantSubscriptionID)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	virtualMachines := NewVirtualMachines(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if virtualMachines.GetSubscriptionID() != wantSubscriptionID {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", virtualMachines.GetSubscriptionID(), wantSubscriptionID)
	}
}
