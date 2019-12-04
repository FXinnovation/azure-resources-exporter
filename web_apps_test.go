package main

import (
	"testing"
)

func TestNewWebApps_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	webApps := NewWebApps(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if webApps.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", webApps.GetSubscriptionID(), want)
	}
}
