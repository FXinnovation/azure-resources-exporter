package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2019-02-01-preview/eventgrid"
	"github.com/prometheus/common/log"
)

const eventGridTopicsResourceType = "Microsoft.EventGrid/topics"

// EventGridTopicsClient is the client implementation to EventGridTopics API
type EventGridTopicsClient struct {
	Session   *AzureSession
	Client    *eventgrid.TopicsClient
	Resources Resources
}

// EventGridTopics client interface
type EventGridTopics interface {
	GetEventGridTopics() (*[]eventgrid.Topic, error)
	GetSubscriptionID() string
}

// NewEventGridTopics returns a new EventGridTopics client
func NewEventGridTopics(session *AzureSession) EventGridTopics {

	client := eventgrid.NewTopicsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &EventGridTopicsClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (ec *EventGridTopicsClient) GetSubscriptionID() string {
	return ec.Session.SubscriptionID
}

// GetEventGridTopics fetch EventGridTopics with state
func (ec *EventGridTopicsClient) GetEventGridTopics() (*[]eventgrid.Topic, error) {
	var egtList []eventgrid.Topic

	resources, err := ec.Resources.GetResources(eventGridTopicsResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping event grid topic: %s", err)
			continue
		}

		egt, err := ec.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		egtList = append(egtList, egt)
	}

	return &egtList, nil
}
