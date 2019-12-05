package main

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// SQLDatabasesCollector collect SQLDatabases metrics
type SQLDatabasesCollector struct {
	sqlDatabases SQLDatabases
}

// NewSQLDatabasesCollector returns the collector
func NewSQLDatabasesCollector(session *AzureSession) *SQLDatabasesCollector {
	sqlDatabases := NewSQLDatabases(session)

	return &SQLDatabasesCollector{
		sqlDatabases: sqlDatabases,
	}
}

// Describe to satisfy the collector interface.
func (c *SQLDatabasesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("SQLDatabasesCollector", "dummy", nil, nil)
}

// Collect metrics from SQL Databases API
func (c *SQLDatabasesCollector) Collect(ch chan<- prometheus.Metric) {

	dbList, err := c.sqlDatabases.GetSQLDatabases()
	if err != nil {
		log.Errorf("Failed to get SQL database list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectDatabaseUp(ch, dbList)
}

// CollectDatabaseUp converts Databases status as a metric
func (c *SQLDatabasesCollector) CollectDatabaseUp(ch chan<- prometheus.Metric, dbList *[]sql.Database) {

	for _, db := range *dbList {
		up := 0.0
		if db.DatabaseProperties.Status == sql.DatabaseStatusOnline {
			up = 1
		}

		labels, err := ParseResourceLabels(*db.ID)
		if err != nil {
			log.Errorf("Skipping SQLDatabase: %s", err)
			continue
		}

		labels["subscription_id"] = c.sqlDatabases.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("sql_database_up", "Status of the SQL database", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
