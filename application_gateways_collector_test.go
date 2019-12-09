package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedApplicationGateways struct {
	mock.Mock
}

func (mock *MockedApplicationGateways) GetApplicationGateways() (*[]network.ApplicationGateway, error) {
	args := mock.Called()
	return args.Get(0).(*[]network.ApplicationGateway), args.Error(1)
}

func (mock *MockedApplicationGateways) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewApplicationGatewaysCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewApplicationGatewaysCollector(session)
}

func TestCollectApplicationGateway_Error(t *testing.T) {
	applicationGateways := MockedApplicationGateways{}
	collector := ApplicationGatewaysCollector{
		applicationGateways: &applicationGateways,
	}
	prometheus.MustRegister(&collector)

	var agList []network.ApplicationGateway
	applicationGateways.On("GetApplicationGateways").Return(&agList, errors.New("Unit test Error"))
	applicationGateways.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectApplicationGateway_Up(t *testing.T) {
	applicationGateways := MockedApplicationGateways{}
	collector := ApplicationGatewaysCollector{
		applicationGateways: &applicationGateways,
	}

	var agList []network.ApplicationGateway
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Network/applicationGateways/my_ag"

	agList = append(agList, network.ApplicationGateway{
		ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
			OperationalState: network.Running,
		},
		ID: &id,
	})

	applicationGateways.On("GetApplicationGateways").Return(&agList, nil)
	applicationGateways.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP application_gateway_up Operational state of the application gateway
# TYPE application_gateway_up gauge
application_gateway_up{resource_group="my_rg",resource_name="my_ag",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
