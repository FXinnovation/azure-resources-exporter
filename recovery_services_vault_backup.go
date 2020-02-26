package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/prometheus/common/log"
)

const recoveryServicesVaultResourceType = "Microsoft.RecoveryServices/vaults"

// RecoveryServicesBackupClient is the client implementation to recoveryservices backup API
type RecoveryServicesBackupClient struct {
	Session                   *AzureSession
	ProtectedItemsGroupClient *backup.ProtectedItemsGroupClient
	Resources                 Resources
}

// RecoveryServicesBackup client interface
type RecoveryServicesBackup interface {
	GetAzureIaaSVMProtectedItem() (*[]backup.AzureIaaSVMProtectedItem, error)
	GetSubscriptionID() string
}

// NewRecoveryServicesBackup returns a new RecoveryServicesBackup client
func NewRecoveryServicesBackup(session *AzureSession) RecoveryServicesBackup {
	client := backup.NewProtectedItemsGroupClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &RecoveryServicesBackupClient{
		Session:                   session,
		ProtectedItemsGroupClient: &client,
		Resources:                 resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (rc *RecoveryServicesBackupClient) GetSubscriptionID() string {
	return rc.Session.SubscriptionID
}

// GetAzureIaaSVMProtectedItem fetch AzureIaaSVMProtectedItems
func (rc *RecoveryServicesBackupClient) GetAzureIaaSVMProtectedItem() (*[]backup.AzureIaaSVMProtectedItem, error) {
	var piList []backup.AzureIaaSVMProtectedItem

	resources, err := rc.Resources.GetResources(recoveryServicesVaultResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceID(*resource.ID)
		if err != nil {
			log.Errorf("Skipping recovery services vault: %s", err)
			continue
		}

		for it, err := rc.ProtectedItemsGroupClient.ListComplete(context.Background(), *resource.Name, labels["resource_group"], "backupManagementType eq 'AzureIaasVM' and itemType eq 'VM'", ""); it.NotDone(); err = it.Next() {
			if err != nil {
				return nil, err
			}
			pir := it.Value()
			pi, _ := pir.Properties.AsAzureIaaSVMProtectedItem()
			piList = append(piList, *pi)
		}
	}

	return &piList, nil
}
