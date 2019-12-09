package main

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// ApplicationGatewaysCollector collect ApplicationGateways metrics
type ApplicationGatewaysCollector struct {
	applicationGateways ApplicationGateways
}

// NewApplicationGatewaysCollector returns the collector
func NewApplicationGatewaysCollector(session *AzureSession) *ApplicationGatewaysCollector {
	applicationGateways := NewApplicationGateways(session)

	return &ApplicationGatewaysCollector{
		applicationGateways: applicationGateways,
	}
}

// Describe to satisfy the collector interface.
func (c *ApplicationGatewaysCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("ApplicationGatewaysCollector", "dummy", nil, nil)
}

// Collect metrics from Application Gateways API
func (c *ApplicationGatewaysCollector) Collect(ch chan<- prometheus.Metric) {

	agList, err := c.applicationGateways.GetApplicationGateways()
	if err != nil {
		log.Errorf("Failed to get application gateway list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectApplicationGatewayUp(ch, agList)
}

// CollectApplicationGatewayUp converts ApplicationGateway operational state as a metric
func (c *ApplicationGatewaysCollector) CollectApplicationGatewayUp(ch chan<- prometheus.Metric, agList *[]network.ApplicationGateway) {

	for _, ag := range *agList {
		up := 0.0
		if ag.ApplicationGatewayPropertiesFormat.OperationalState == network.Running {
			up = 1
		}

		labels, err := ParseResourceID(*ag.ID)
		if err != nil {
			log.Errorf("Skipping Application Gateway: %s", err)
			continue
		}

		labels["subscription_id"] = c.applicationGateways.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("application_gateway_up", "Operational state of the application gateway", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
