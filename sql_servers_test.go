package main

import (
	"testing"
)

func TestNewSQLServers_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	sqlServers := NewSQLServers(session)

	if sqlServers.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", sqlServers.GetSubscriptionID(), want)
	}
}
