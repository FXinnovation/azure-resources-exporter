package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedAutomationAccounts struct {
	mock.Mock
}

func (mock *MockedAutomationAccounts) GetAutomationAccounts() (*[]automation.Account, error) {
	args := mock.Called()
	return args.Get(0).(*[]automation.Account), args.Error(1)
}

func (mock *MockedAutomationAccounts) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewAutomationAccountsCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewAutomationAccountsCollector(session)
}

func TestCollectAutomationAccount_Error(t *testing.T) {
	automationAccounts := MockedAutomationAccounts{}
	collector := AutomationAccountsCollector{
		automationAccounts: &automationAccounts,
	}
	prometheus.MustRegister(&collector)

	var conList []automation.Account
	automationAccounts.On("GetAutomationAccounts").Return(&conList, errors.New("Unit test Error"))
	automationAccounts.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectAutomationAccount_Up(t *testing.T) {
	automationAccounts := MockedAutomationAccounts{}
	collector := AutomationAccountsCollector{
		automationAccounts: &automationAccounts,
	}

	var agList []automation.Account
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Automation/automationAccounts/my_aa"

	agList = append(agList, automation.Account{
		AccountProperties: &automation.AccountProperties{
			State: automation.Ok,
		},
		ID: &id,
	})

	automationAccounts.On("GetAutomationAccounts").Return(&agList, nil)
	automationAccounts.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP automation_account_up State of the automation account
# TYPE automation_account_up gauge
automation_account_up{resource_group="my_rg",resource_name="my_aa",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
