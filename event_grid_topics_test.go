package main

import (
	"testing"
)

func TestNewEventGridTopics_OK(t *testing.T) {
	want := "subscriptionID"
	session, err := NewAzureSession(want)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	eventGridTopics := NewEventGridTopics(session)

	if eventGridTopics.GetSubscriptionID() != want {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", eventGridTopics.GetSubscriptionID(), want)
	}
}
