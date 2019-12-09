package main

import (
	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// AutomationAccountsCollector collect AutomationAccounts metrics
type AutomationAccountsCollector struct {
	automationAccounts AutomationAccounts
}

// NewAutomationAccountsCollector returns the collector
func NewAutomationAccountsCollector(session *AzureSession) *AutomationAccountsCollector {
	automationAccounts := NewAutomationAccounts(session)

	return &AutomationAccountsCollector{
		automationAccounts: automationAccounts,
	}
}

// Describe to satisfy the collector interface.
func (c *AutomationAccountsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("AutomationAccountsCollector", "dummy", nil, nil)
}

// Collect metrics from AutomationAccounts API
func (c *AutomationAccountsCollector) Collect(ch chan<- prometheus.Metric) {

	aaList, err := c.automationAccounts.GetAutomationAccounts()
	if err != nil {
		log.Errorf("Failed to get automation account list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectAutomationAccountUp(ch, aaList)
}

// CollectAutomationAccountUp converts AutomationAccount state as a metric
func (c *AutomationAccountsCollector) CollectAutomationAccountUp(ch chan<- prometheus.Metric, aaList *[]automation.Account) {

	for _, aa := range *aaList {
		up := 0.0
		if aa.AccountProperties.State == automation.Ok {
			up = 1
		}

		labels, err := ParseResourceID(*aa.ID)
		if err != nil {
			log.Errorf("Skipping Automation Account: %s", err)
			continue
		}

		labels["subscription_id"] = c.automationAccounts.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("automation_account_up", "State of the automation account", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
