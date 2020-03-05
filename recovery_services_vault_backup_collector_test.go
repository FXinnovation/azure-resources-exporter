package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedRecoveryServicesBackup struct {
	mock.Mock
}

func (mock *MockedRecoveryServicesBackup) GetAzureIaaSComputeVMProtectedItem() (*[]backup.AzureIaaSComputeVMProtectedItem, error) {
	args := mock.Called()
	return args.Get(0).(*[]backup.AzureIaaSComputeVMProtectedItem), args.Error(1)
}

func (mock *MockedRecoveryServicesBackup) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewRecoveryServicesBackupCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewRecoveryServicesBackupCollector(session)
}

func TestCollectAzureIaaSComputeVMProtectedItem_Error(t *testing.T) {
	recoveryServicesBackup := MockedRecoveryServicesBackup{}
	collector := RecoveryServicesBackupCollector{
		recoveryServicesBackup: &recoveryServicesBackup,
	}
	prometheus.MustRegister(&collector)

	var piList []backup.AzureIaaSComputeVMProtectedItem
	recoveryServicesBackup.On("GetAzureIaaSComputeVMProtectedItem").Return(&piList, errors.New("Unit test Error"))
	recoveryServicesBackup.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectAzureIaaSComputeVMProtectedItem_Up(t *testing.T) {
	recoveryServicesBackup := MockedRecoveryServicesBackup{}
	collector := RecoveryServicesBackupCollector{
		recoveryServicesBackup: &recoveryServicesBackup,
	}

	var piList []backup.AzureIaaSComputeVMProtectedItem
	virtualMachineID := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Compute/virtualMachines/my_vm"
	lastBackupStatus := string(backup.JobStatusCompleted)

	lastBackupTime := date.Time{}
	lastBackupTime.UnmarshalText([]byte("2012-11-24T00:00:00Z"))

	piList = append(piList, backup.AzureIaaSComputeVMProtectedItem{
		VirtualMachineID: &virtualMachineID,
		HealthStatus:     backup.HealthStatusPassed,
		LastBackupStatus: &lastBackupStatus,
		LastBackupTime:   &lastBackupTime,
	})

	recoveryServicesBackup.On("GetAzureIaaSComputeVMProtectedItem").Return(&piList, nil)
	recoveryServicesBackup.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP recovery_services_vault_backup_item_vm_health_status Health status of the VM backup
# TYPE recovery_services_vault_backup_item_vm_health_status gauge
recovery_services_vault_backup_item_vm_health_status{health_status="ActionRequired",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_health_status{health_status="ActionSuggested",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_health_status{health_status="Invalid",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_health_status{health_status="Passed",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 1
# HELP recovery_services_vault_backup_item_vm_last_backup_status Last backup status of the VM backup
# TYPE recovery_services_vault_backup_item_vm_last_backup_status gauge
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="Cancelled",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="Cancelling",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="Completed",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 1
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="CompletedWithWarnings",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="Failed",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="InProgress",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
recovery_services_vault_backup_item_vm_last_backup_status{last_backup_status="Invalid",resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
# HELP recovery_services_vault_backup_item_vm_last_backup_time_seconds Unix/epoch time of the last VM backup
# TYPE recovery_services_vault_backup_item_vm_last_backup_time_seconds gauge
recovery_services_vault_backup_item_vm_last_backup_time_seconds{resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 1.3537152e+09
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
