package main

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// ServerUpState is the state of a Server that is UP
const ServerUpState = "Ready"

// SQLServersCollector collect SQLServers metrics
type SQLServersCollector struct {
	sqlServers SQLServers
}

// NewSQLServersCollector returns the collector
func NewSQLServersCollector(session *AzureSession) *SQLServersCollector {
	sqlServers := NewSQLServers(session)

	return &SQLServersCollector{
		sqlServers: sqlServers,
	}
}

// Describe to satisfy the collector interface.
func (c *SQLServersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("SQLServersCollector", "dummy", nil, nil)
}

// Collect metrics from SQL Server API
func (c *SQLServersCollector) Collect(ch chan<- prometheus.Metric) {

	serverList, err := c.sqlServers.GetSQLServers()
	if err != nil {
		log.Errorf("Failed to get SQL server list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectServerUp(ch, serverList)
}

// CollectServerUp converts SQLServers status as a metric
func (c *SQLServersCollector) CollectServerUp(ch chan<- prometheus.Metric, serverList *[]sql.Server) {

	for _, server := range *serverList {
		up := 0.0
		if *server.ServerProperties.State == ServerUpState {
			up = 1
		}

		labels, err := ParseResourceID(*server.ID)
		if err != nil {
			log.Errorf("Skipping SQLServer: %s", err)
			continue
		}

		labels["subscription_id"] = c.sqlServers.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("sql_server_up", "Status of the SQL server", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
