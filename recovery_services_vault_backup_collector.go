package main

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// RecoveryServicesBackupCollector collect RecoveryServicesBackup metrics
type RecoveryServicesBackupCollector struct {
	recoveryServicesBackup RecoveryServicesBackup
}

// NewRecoveryServicesBackupCollector returns the collector
func NewRecoveryServicesBackupCollector(session *AzureSession) *RecoveryServicesBackupCollector {
	recoveryServicesBackup := NewRecoveryServicesBackup(session)

	return &RecoveryServicesBackupCollector{
		recoveryServicesBackup: recoveryServicesBackup,
	}
}

// Describe to satisfy the collector interface.
func (c *RecoveryServicesBackupCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("RecoveryServicesBackupCollector", "dummy", nil, nil)
}

// Collect metrics from RecoveryServicesBackup API
func (c *RecoveryServicesBackupCollector) Collect(ch chan<- prometheus.Metric) {

	vmProtectedItems, err := c.recoveryServicesBackup.GetAzureIaaSComputeVMProtectedItem()
	if err != nil {
		log.Errorf("Failed to get AzureIaaSVM Protected Item list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectAzureIaaSComputeVMProtectedItem(ch, vmProtectedItems)
}

// CollectAzureIaaSComputeVMProtectedItem converts AzureIaaSComputeVM Protected Items status as a metric
func (c *RecoveryServicesBackupCollector) CollectAzureIaaSComputeVMProtectedItem(ch chan<- prometheus.Metric, vmProtectedItems *[]backup.AzureIaaSComputeVMProtectedItem) {

	possibleHealthStatuses := backup.PossibleHealthStatusValues()
	possibleLastBackupStatuses := backup.PossibleJobStatusValues()

	for _, vmProtectedItem := range *vmProtectedItems {
		labels, err := ParseResourceID(*vmProtectedItem.VirtualMachineID)
		if err != nil {
			log.Errorf("Skipping VM protected item: %s", err)
			continue
		}
		labels["subscription_id"] = c.recoveryServicesBackup.GetSubscriptionID()

		healthStatus := vmProtectedItem.HealthStatus
		healthLabel := make(map[string]string)
		for key, value := range labels {
			healthLabel[key] = value
		}
		for _, possibleHealthStatus := range possibleHealthStatuses {
			status := 0.0
			if possibleHealthStatus == healthStatus {
				status = 1
			}
			healthLabel["health_status"] = string(possibleHealthStatus)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("recovery_services_vault_backup_item_vm_health_status", "Health status of the VM backup", nil, healthLabel),
				prometheus.GaugeValue,
				status,
			)
		}

		lastBackupStatus := vmProtectedItem.LastBackupStatus
		lastBackupLabel := make(map[string]string)
		for key, value := range labels {
			lastBackupLabel[key] = value
		}
		for _, possibleLastBackupStatus := range possibleLastBackupStatuses {
			status := 0.0
			if string(possibleLastBackupStatus) == *lastBackupStatus {
				status = 1
			}
			lastBackupLabel["last_backup_status"] = string(possibleLastBackupStatus)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("recovery_services_vault_backup_item_vm_last_backup_status", "Last backup status of the VM backup", nil, lastBackupLabel),
				prometheus.GaugeValue,
				status,
			)
		}
	}
}
