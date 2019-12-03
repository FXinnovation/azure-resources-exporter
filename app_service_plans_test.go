package main

import (
	"testing"
)

func TestNewAppServicePlans_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	appServicePlans := NewAppServicePlans(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if appServicePlans.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", appServicePlans.GetSubscriptionID(), want)
	}
}
