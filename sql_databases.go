package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/prometheus/common/log"
)

// SQLDatabasesClient is the client implementation to SQLDatabases API
type SQLDatabasesClient struct {
	Session   *AzureSession
	Client    *sql.DatabasesClient
	Resources Resources
}

// SQLDatabases client interface
type SQLDatabases interface {
	GetSQLDatabases() (*[]sql.Database, error)
	GetSubscriptionID() string
}

// NewSQLDatabases returns a new SQL Databases client
func NewSQLDatabases(session *AzureSession) SQLDatabases {
	client := sql.NewDatabasesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &SQLDatabasesClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (sc *SQLDatabasesClient) GetSubscriptionID() string {
	return sc.Session.SubscriptionID
}

// GetSQLDatabases fetch SQL Databases with status
func (sc *SQLDatabasesClient) GetSQLDatabases() (*[]sql.Database, error) {
	var dbList []sql.Database

	resources, err := sc.Resources.GetResources(sqlServersResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping SQL server: %s", err)
			continue
		}

		for it, err := sc.Client.ListByServerComplete(context.Background(), labels["resource_group"], *resource.Name); it.NotDone(); err = it.Next() {
			if err != nil {
				return nil, err
			}
			dbList = append(dbList, it.Value())
		}
	}

	return &dbList, nil
}
