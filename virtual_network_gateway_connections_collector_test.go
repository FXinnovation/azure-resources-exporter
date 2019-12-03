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

type MockedVirtualNetworkGatewayConnections struct {
	mock.Mock
}

func (mock *MockedVirtualNetworkGatewayConnections) GetVirtualNetworkGatewayConnections() (*[]network.VirtualNetworkGatewayConnection, error) {
	args := mock.Called()
	return args.Get(0).(*[]network.VirtualNetworkGatewayConnection), args.Error(1)
}

func (mock *MockedVirtualNetworkGatewayConnections) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewVirtualNetworkGatewayConnecitonsCollector_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	session, err := NewAzureSession(wantSubscriptionID)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewVirtualNetworkGatewayConnectionsCollector(session)

	if err != nil {
		t.Errorf("Error occured %s", err)
	}
}

func TestCollectVNGCon_Error(t *testing.T) {
	virtualNetworkGatewayConnections := MockedVirtualNetworkGatewayConnections{}
	collector := VirtualNetworkGatewayConnectionsCollector{
		virtualNetworkGatewayConnections: &virtualNetworkGatewayConnections,
	}
	prometheus.MustRegister(&collector)

	var conList []network.VirtualNetworkGatewayConnection
	virtualNetworkGatewayConnections.On("GetVirtualNetworkGatewayConnections").Return(&conList, errors.New("Unit test Error"))
	virtualNetworkGatewayConnections.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectVNGCon_Up(t *testing.T) {
	virtualNetworkGatewayConnections := MockedVirtualNetworkGatewayConnections{}
	collector := VirtualNetworkGatewayConnectionsCollector{
		virtualNetworkGatewayConnections: &virtualNetworkGatewayConnections,
	}

	var conList []network.VirtualNetworkGatewayConnection
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Network/connections/my_con"

	conList = append(conList, network.VirtualNetworkGatewayConnection{
		VirtualNetworkGatewayConnectionPropertiesFormat: &network.VirtualNetworkGatewayConnectionPropertiesFormat{
			ConnectionStatus: network.VirtualNetworkGatewayConnectionStatusConnected,
		},
		ID: &id,
	})

	virtualNetworkGatewayConnections.On("GetVirtualNetworkGatewayConnections").Return(&conList, nil)
	virtualNetworkGatewayConnections.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP virtual_network_gateway_connection_up Connections status of the virtual network gateway
# TYPE virtual_network_gateway_connection_up gauge
virtual_network_gateway_connection_up{resource_group="my_rg",resource_name="my_con",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
