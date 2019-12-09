package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/prometheus/common/log"
)

const automationAccountsResourceType = "Microsoft.Automation/automationAccounts"

// AutomationAccountsClient is the client implementation to AutomationAccounts API
type AutomationAccountsClient struct {
	Session   *AzureSession
	Client    *automation.AccountClient
	Resources Resources
}

// AutomationAccounts client interface
type AutomationAccounts interface {
	GetAutomationAccounts() (*[]automation.Account, error)
	GetSubscriptionID() string
}

// NewAutomationAccounts returns a new AutomationAccounts client
func NewAutomationAccounts(session *AzureSession) AutomationAccounts {

	client := automation.NewAccountClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &AutomationAccountsClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (ac *AutomationAccountsClient) GetSubscriptionID() string {
	return ac.Session.SubscriptionID
}

// GetAutomationAccounts fetch AutomationAccounts with state
func (ac *AutomationAccountsClient) GetAutomationAccounts() (*[]automation.Account, error) {
	var aaList []automation.Account

	resources, err := ac.Resources.GetResources(automationAccountsResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping automation account: %s", err)
			continue
		}

		aa, err := ac.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		aaList = append(aaList, aa)
	}

	return &aaList, nil
}
