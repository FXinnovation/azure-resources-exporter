package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedServiceFabricClusters struct {
	mock.Mock
}

func (mock *MockedServiceFabricClusters) GetServiceFabricClusters() (*[]servicefabric.Cluster, error) {
	args := mock.Called()
	return args.Get(0).(*[]servicefabric.Cluster), args.Error(1)
}

func (mock *MockedServiceFabricClusters) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewServiceFabricClustersCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewServiceFabricClustersCollector(session)
}

func TestCollectServiceFabricCluster_Error(t *testing.T) {
	serviceFabricClusters := MockedServiceFabricClusters{}
	collector := ServiceFabricClustersCollector{
		serviceFabricClusters: &serviceFabricClusters,
	}
	prometheus.MustRegister(&collector)

	var sfcList []servicefabric.Cluster
	serviceFabricClusters.On("GetServiceFabricClusters").Return(&sfcList, errors.New("Unit test Error"))
	serviceFabricClusters.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectServiceFabricCluster_Up(t *testing.T) {
	serviceFabricClusters := MockedServiceFabricClusters{}
	collector := ServiceFabricClustersCollector{
		serviceFabricClusters: &serviceFabricClusters,
	}

	var sfcList []servicefabric.Cluster
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.ServiceFabric/clusters/my_sfc"
	tagValue := "Value"
	resourceType := "Microsoft.ServiceFabric/clusters"

	sfcList = append(sfcList, servicefabric.Cluster{
		ClusterProperties: &servicefabric.ClusterProperties{
			ClusterState: servicefabric.Ready,
		},
		ID: &id,
		Tags: map[string]*string{
			"Key": &tagValue,
		},
		Type: &resourceType,
	})

	serviceFabricClusters.On("GetServiceFabricClusters").Return(&sfcList, nil)
	serviceFabricClusters.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP azure_tag_info Tags of the Azure resource
# TYPE azure_tag_info gauge
azure_tag_info{resource_group="my_rg",resource_name="my_sfc",resource_type="Microsoft.ServiceFabric/clusters",subscription_id="my_subscription",tag_key="Value"} 1
# HELP service_fabric_cluster_up State of the service fabric cluster
# TYPE service_fabric_cluster_up gauge
service_fabric_cluster_up{resource_group="my_rg",resource_name="my_sfc",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
