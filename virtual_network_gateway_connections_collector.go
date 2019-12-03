package main

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// VirtualNetworkGatewayConnecitonsCollector collect VirtualNetworkGatewayConnecitons metrics
type VirtualNetworkGatewayConnecitonsCollector struct {
	virtualNetworkGatewayConnections VirtualNetworkGatewayConnections
}

// NewVirtualNetworkGatewayConnectionsCollector returns the collector
func NewVirtualNetworkGatewayConnectionsCollector(session *AzureSession) *VirtualNetworkGatewayConnecitonsCollector {
	virtualNetworkGatewayConnections := NewVirtualNetworkGatewayConnections(session)

	return &VirtualNetworkGatewayConnecitonsCollector{
		virtualNetworkGatewayConnections: virtualNetworkGatewayConnections,
	}
}

// Describe to satisfy the collector interface.
func (c *VirtualNetworkGatewayConnecitonsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("VirtualNetworkGatewayConnecitonsCollector", "dummy", nil, nil)
}

// Collect metrics from Virtual Machine API
func (c *VirtualNetworkGatewayConnecitonsCollector) Collect(ch chan<- prometheus.Metric) {

	conList, err := c.virtualNetworkGatewayConnections.GetVirtualNetworkGatewayConnections()

	if err != nil {
		log.Errorf("Failed to get virtual network gateway connections list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectConnectionUp(ch, conList)
}

// CollectConnectionUp converts VirtualNetworkGateway connections status as a metric
func (c *VirtualNetworkGatewayConnecitonsCollector) CollectConnectionUp(ch chan<- prometheus.Metric, conList *[]network.VirtualNetworkGatewayConnection) {

	for _, con := range *conList {
		up := 0.0
		if con.VirtualNetworkGatewayConnectionPropertiesFormat.ConnectionStatus == network.VirtualNetworkGatewayConnectionStatusConnected {
			up = 1
		}

		labels, err := ParseResourceLabels(*con.ID)

		if err != nil {
			log.Errorf("Skipping VirtualNetworkGateway connection: %s", err)
			continue
		}

		labels["subscription_id"] = c.virtualNetworkGatewayConnections.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("virtual_network_gateway_connection_up", "Connections status of the virtual network gateway", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
