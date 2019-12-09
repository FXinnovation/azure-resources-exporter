package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/prometheus/common/log"
)

const sqlServersResourceType = "Microsoft.Sql/servers"

// SQLServersClient is the client implementation to SQLServers API
type SQLServersClient struct {
	Session   *AzureSession
	Client    *sql.ServersClient
	Resources Resources
}

// SQLServers client interface
type SQLServers interface {
	GetSQLServers() (*[]sql.Server, error)
	GetSubscriptionID() string
}

// NewSQLServers returns a new SQL Servers client
func NewSQLServers(session *AzureSession) SQLServers {
	client := sql.NewServersClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &SQLServersClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (sc *SQLServersClient) GetSubscriptionID() string {
	return sc.Session.SubscriptionID
}

// GetSQLServers fetch SQL Servers with state
func (sc *SQLServersClient) GetSQLServers() (*[]sql.Server, error) {
	var serverList []sql.Server

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

		server, err := sc.Client.Get(context.Background(), labels["resource_group"], *resource.Name)
		if err != nil {
			return nil, err
		}

		serverList = append(serverList, server)
	}

	return &serverList, nil
}
