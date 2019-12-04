package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedSQLServers struct {
	mock.Mock
}

func (mock *MockedSQLServers) GetSQLServers() (*[]sql.Server, error) {
	args := mock.Called()
	return args.Get(0).(*[]sql.Server), args.Error(1)
}

func (mock *MockedSQLServers) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewSQLServersCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewSQLServersCollector(session)
}

func TestCollectServer_Error(t *testing.T) {
	sqlServers := MockedSQLServers{}
	collector := SQLServersCollector{
		sqlServers: &sqlServers,
	}
	prometheus.MustRegister(&collector)

	var serverList []sql.Server
	sqlServers.On("GetSQLServers").Return(&serverList, errors.New("Unit test Error"))
	sqlServers.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectServer_Up(t *testing.T) {
	sqlServers := MockedSQLServers{}
	collector := SQLServersCollector{
		sqlServers: &sqlServers,
	}

	var serverList []sql.Server
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Sql/servers/my_server"
	state := "Ready"

	serverList = append(serverList, sql.Server{
		ServerProperties: &sql.ServerProperties{
			State: &state,
		},
		ID: &id,
	})

	sqlServers.On("GetSQLServers").Return(&serverList, nil)
	sqlServers.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP sql_server_up Status of the SQL server
# TYPE sql_server_up gauge
sql_server_up{resource_group="my_rg",resource_name="my_server",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
