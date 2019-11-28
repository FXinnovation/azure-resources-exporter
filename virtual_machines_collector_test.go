package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedVirtualMachines struct {
	mock.Mock
}

func (mock *MockedVirtualMachines) GetVirtualMachines() (*[]compute.VirtualMachine, error) {
	args := mock.Called()
	return args.Get(0).(*[]compute.VirtualMachine), args.Error(1)
}

func (mock *MockedVirtualMachines) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewVirtualMachinesCollector_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	session, err := NewAzureSession(wantSubscriptionID)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewVirtualMachinesCollector(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}
}

func TestCollect_Error(t *testing.T) {
	virtualMachines := MockedVirtualMachines{}
	collector := VirtualMachinesCollector{
		virtualMachines: &virtualMachines,
	}
	prometheus.MustRegister(&collector)

	var vmList []compute.VirtualMachine
	virtualMachines.On("GetVirtualMachines").Return(&vmList, errors.New("Unit test Error"))
	virtualMachines.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollect_No_Status(t *testing.T) {
	virtualMachines := MockedVirtualMachines{}
	collector := VirtualMachinesCollector{
		virtualMachines: &virtualMachines,
	}

	var vmList []compute.VirtualMachine
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Compute/virtualMachines/my_vm"
	vmList = append(vmList, compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			InstanceView: &compute.VirtualMachineInstanceView{
				Statuses: &[]compute.InstanceViewStatus{},
			},
		},
		ID: &id,
	})

	virtualMachines.On("GetVirtualMachines").Return(&vmList, nil)
	virtualMachines.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP virtual_machine_instance_up Running status of the virtual machine instance
# TYPE virtual_machine_instance_up gauge
virtual_machine_instance_up{resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 0
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}

func TestCollect_Up(t *testing.T) {
	virtualMachines := MockedVirtualMachines{}
	collector := VirtualMachinesCollector{
		virtualMachines: &virtualMachines,
	}

	var vmList []compute.VirtualMachine
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Compute/virtualMachines/my_vm"
	runningStatus := VirtualMachineUpState
	otherStatus := "other"
	vmList = append(vmList, compute.VirtualMachine{
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			InstanceView: &compute.VirtualMachineInstanceView{
				Statuses: &[]compute.InstanceViewStatus{
					compute.InstanceViewStatus{
						Code: &runningStatus,
					},
					compute.InstanceViewStatus{
						Code: &otherStatus,
					},
				},
			},
		},
		ID: &id,
	})

	virtualMachines.On("GetVirtualMachines").Return(&vmList, nil)
	virtualMachines.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP virtual_machine_instance_up Running status of the virtual machine instance
# TYPE virtual_machine_instance_up gauge
virtual_machine_instance_up{resource_group="my_rg",resource_name="my_vm",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
