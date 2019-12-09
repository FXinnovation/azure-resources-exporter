package main

import (
	"testing"
)

func TestNewAutomationAccounts_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	automationAccounts := NewAutomationAccounts(session)

	if automationAccounts.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", automationAccounts.GetSubscriptionID(), want)
	}
}
