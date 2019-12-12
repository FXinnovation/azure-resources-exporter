package main

import (
	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// ServiceFabricClustersCollector collect ServiceFabricClusters metrics
type ServiceFabricClustersCollector struct {
	serviceFabricClusters ServiceFabricClusters
}

// NewServiceFabricClustersCollector returns the collector
func NewServiceFabricClustersCollector(session *AzureSession) *ServiceFabricClustersCollector {
	serviceFabricClusters := NewServiceFabricClusters(session)

	return &ServiceFabricClustersCollector{
		serviceFabricClusters: serviceFabricClusters,
	}
}

// Describe to satisfy the collector interface.
func (c *ServiceFabricClustersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("ServiceFabricClustersCollector", "dummy", nil, nil)
}

// Collect metrics from ServiceFabricClusters API
func (c *ServiceFabricClustersCollector) Collect(ch chan<- prometheus.Metric) {

	sfcList, err := c.serviceFabricClusters.GetServiceFabricClusters()
	if err != nil {
		log.Errorf("Failed to get Service fabric cluster list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectServiceFabricClusterUp(ch, sfcList)
}

// CollectServiceFabricClusterUp converts ServiceFabricCluster state as a metric
func (c *ServiceFabricClustersCollector) CollectServiceFabricClusterUp(ch chan<- prometheus.Metric, sfcList *[]servicefabric.Cluster) {

	for _, sfc := range *sfcList {
		up := 0.0
		if sfc.ClusterProperties.ClusterState == servicefabric.Ready {
			up = 1
		}

		labels, err := ParseResourceID(*sfc.ID)
		if err != nil {
			log.Errorf("Skipping Service Fabric Cluster: %s", err)
			continue
		}

		labels["subscription_id"] = c.serviceFabricClusters.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("service_fabric_cluster_up", "State of the service fabric cluster", nil, labels),
			prometheus.GaugeValue,
			up,
		)

		ExportAzureTagInfo(ch, sfc.Tags, sfc.Type, labels)
	}
}
