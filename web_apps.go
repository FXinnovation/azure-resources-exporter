package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/prometheus/common/log"
)

const webAppsResourceType = "Microsoft.Web/sites"

// WebAppsClient is the client implementation to WebApps API
type WebAppsClient struct {
	Session   *AzureSession
	Client    *web.AppsClient
	Resources Resources
}

// WebApps client interface
type WebApps interface {
	GetWebApps() (*[]web.Site, error)
	GetSubscriptionID() string
}

// NewWebApps returns a new Web App client
func NewWebApps(session *AzureSession) WebApps {
	client := web.NewAppsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &WebAppsClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (wc *WebAppsClient) GetSubscriptionID() string {
	return wc.Session.SubscriptionID
}

// GetWebApps fetch Web Apps with state
func (wc *WebAppsClient) GetWebApps() (*[]web.Site, error) {
	var appList []web.Site

	resources, err := wc.Resources.GetResources(webAppsResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceLabels(*resource.ID)
		if err != nil {
			log.Errorf("Skipping web app: %s", err)
			continue
		}

		app, err := wc.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		appList = append(appList, app)
	}

	return &appList, nil
}
