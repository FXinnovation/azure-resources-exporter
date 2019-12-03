package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/prometheus/common/log"
)

const appServicePlansResourceType = "Microsoft.Web/serverfarms"

// AppServicePlansClient is the client implementation to AppServicePlans API
type AppServicePlansClient struct {
	Session   *AzureSession
	Client    *web.AppServicePlansClient
	Resources Resources
}

// AppServicePlans client interface
type AppServicePlans interface {
	GetAppServicePlans() (*[]web.AppServicePlan, error)
	GetSubscriptionID() string
}

// NewAppServicePlans returns a new App Service Plans client
func NewAppServicePlans(session *AzureSession) AppServicePlans {
	client := web.NewAppServicePlansClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &AppServicePlansClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (ac *AppServicePlansClient) GetSubscriptionID() string {
	return ac.Session.SubscriptionID
}

// GetAppServicePlans fetch App Service Plans with status
func (ac *AppServicePlansClient) GetAppServicePlans() (*[]web.AppServicePlan, error) {
	var planList []web.AppServicePlan

	ressources, err := ac.Resources.GetResources(appServicePlansResourceType)
	if err != nil {
		return nil, err
	}

	for _, ressource := range *ressources {
		labels, err := ParseResourceLabels(*ressource.ID)
		if err != nil {
			log.Errorf("Skipping app service plan: %s", err)
			continue
		}

		plan, err := ac.Client.Get(context.Background(), labels["resource_group"], *ressource.Name)
		if err != nil {
			return nil, err
		}

		planList = append(planList, plan)
	}

	return &planList, nil
}
