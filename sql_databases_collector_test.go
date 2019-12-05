package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedSQLDatabases struct {
	mock.Mock
}

func (mock *MockedSQLDatabases) GetSQLDatabases() (*[]sql.Database, error) {
	args := mock.Called()
	return args.Get(0).(*[]sql.Database), args.Error(1)
}

func (mock *MockedSQLDatabases) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewSQLDatabasesCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewSQLDatabasesCollector(session)
}

func TestCollectDatabase_Error(t *testing.T) {
	sqlDatabases := MockedSQLDatabases{}
	collector := SQLDatabasesCollector{
		sqlDatabases: &sqlDatabases,
	}
	prometheus.MustRegister(&collector)

	var dbList []sql.Database
	sqlDatabases.On("GetSQLDatabases").Return(&dbList, errors.New("Unit test Error"))
	sqlDatabases.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectDatabase_Up(t *testing.T) {
	sqlDatabases := MockedSQLDatabases{}
	collector := SQLDatabasesCollector{
		sqlDatabases: &sqlDatabases,
	}

	var dbList []sql.Database
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Sql/servers/my_server/databases/my_db"

	dbList = append(dbList, sql.Database{
		DatabaseProperties: &sql.DatabaseProperties{
			Status: sql.DatabaseStatusOnline,
		},
		ID: &id,
	})

	sqlDatabases.On("GetSQLDatabases").Return(&dbList, nil)
	sqlDatabases.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP sql_database_up Status of the SQL database
# TYPE sql_database_up gauge
sql_database_up{resource_group="my_rg",resource_name="my_server",sub_resource_name="my_db",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
