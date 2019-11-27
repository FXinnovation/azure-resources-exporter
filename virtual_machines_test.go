package main

import (
	"testing"
)

func TestNewVirtualMachines_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	virtualMachines, err := NewVirtualMachines(wantSubscriptionID)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if virtualMachines.GetSubscriptionID() != wantSubscriptionID {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", virtualMachines.GetSubscriptionID(), wantSubscriptionID)
	}
}

func TestNewVirtualMachines_MissingSubscriptionID(t *testing.T) {
	_, err := NewVirtualMachines("")

	if err == nil {
		t.Errorf("Want an error, got none")
	}
}
