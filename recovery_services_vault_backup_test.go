package main

import (
	"testing"
)

func TestNewRecoveryServicesBackup_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	rsb := NewRecoveryServicesBackup(session)

	if rsb.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", rsb.GetSubscriptionID(), want)
	}
}
