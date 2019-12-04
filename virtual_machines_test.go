package main

import (
	"testing"
)

func TestNewVirtualMachines_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	virtualMachines := NewVirtualMachines(session)

	if virtualMachines.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", virtualMachines.GetSubscriptionID(), want)
	}
}
