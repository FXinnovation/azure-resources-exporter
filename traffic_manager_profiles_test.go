package main

import (
	"testing"
)

func TestNewTrafficManagerProfiles_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	trafficManagerProfiles := NewTrafficManagerProfiles(session)

	if trafficManagerProfiles.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", trafficManagerProfiles.GetSubscriptionID(), want)
	}
}
