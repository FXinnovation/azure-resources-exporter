package main

import (
	"testing"
)

func TestNewSQLDatabases_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	sqlServers := NewSQLDatabases(session)

	if sqlServers.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", sqlServers.GetSubscriptionID(), want)
	}
}
