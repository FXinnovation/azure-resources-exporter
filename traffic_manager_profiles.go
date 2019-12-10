package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/prometheus/common/log"
)

const trafficManagerProfilesResourceType = "Microsoft.Network/trafficmanagerprofiles"

// TrafficManagerProfilesClient is the client implementation to TrafficManagerProfiles API
type TrafficManagerProfilesClient struct {
	Session   *AzureSession
	Client    *trafficmanager.ProfilesClient
	Resources Resources
}

// TrafficManagerProfiles client interface
type TrafficManagerProfiles interface {
	GetTrafficManagerProfiles() (*[]trafficmanager.Profile, error)
	GetSubscriptionID() string
}

// NewTrafficManagerProfiles returns a new TrafficManagerProfiles client
func NewTrafficManagerProfiles(session *AzureSession) TrafficManagerProfiles {

	client := trafficmanager.NewProfilesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &TrafficManagerProfilesClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (tc *TrafficManagerProfilesClient) GetSubscriptionID() string {
	return tc.Session.SubscriptionID
}

// GetTrafficManagerProfiles fetch TrafficManagerProfiles with state
func (tc *TrafficManagerProfilesClient) GetTrafficManagerProfiles() (*[]trafficmanager.Profile, error) {
	var profileList []trafficmanager.Profile

	resources, err := tc.Resources.GetResources(trafficManagerProfilesResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping traffic manager profile: %s", err)
			continue
		}

		profile, err := tc.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		profileList = append(profileList, profile)
	}

	return &profileList, nil
}
