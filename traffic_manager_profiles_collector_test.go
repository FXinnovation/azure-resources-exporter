package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedTrafficManagerProfiles struct {
	mock.Mock
}

func (mock *MockedTrafficManagerProfiles) GetTrafficManagerProfiles() (*[]trafficmanager.Profile, error) {
	args := mock.Called()
	return args.Get(0).(*[]trafficmanager.Profile), args.Error(1)
}

func (mock *MockedTrafficManagerProfiles) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewTrafficManagerProfilesCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewTrafficManagerProfilesCollector(session)
}

func TestCollectTrafficManagerProfile_Error(t *testing.T) {
	trafficManagerProfiles := MockedTrafficManagerProfiles{}
	collector := TrafficManagerProfilesCollector{
		trafficManagerProfiles: &trafficManagerProfiles,
	}
	prometheus.MustRegister(&collector)

	var profileList []trafficmanager.Profile
	trafficManagerProfiles.On("GetTrafficManagerProfiles").Return(&profileList, errors.New("Unit test Error"))
	trafficManagerProfiles.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectTrafficManagerProfile_Up(t *testing.T) {
	trafficManagerProfiles := MockedTrafficManagerProfiles{}
	collector := TrafficManagerProfilesCollector{
		trafficManagerProfiles: &trafficManagerProfiles,
	}

	var profileList []trafficmanager.Profile
	profileID := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Network/trafficManagerProfiles/my_profile"
	endpointID := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Network/trafficManagerProfiles/my_profile/azureEndpoints/my_endpoint"

	profileList = append(profileList, trafficmanager.Profile{
		ProfileProperties: &trafficmanager.ProfileProperties{
			ProfileStatus: trafficmanager.ProfileStatusEnabled,
			Endpoints: &[]trafficmanager.Endpoint{
				trafficmanager.Endpoint{
					EndpointProperties: &trafficmanager.EndpointProperties{
						EndpointStatus: trafficmanager.EndpointStatusEnabled,
					},
					ID: &endpointID,
				},
			},
		},
		ID: &profileID,
	})

	trafficManagerProfiles.On("GetTrafficManagerProfiles").Return(&profileList, nil)
	trafficManagerProfiles.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP traffic_manager_profile_endpoint_up Status of the traffic manager profile endpoint
# TYPE traffic_manager_profile_endpoint_up gauge
traffic_manager_profile_endpoint_up{resource_group="my_rg",resource_name="my_profile",sub_resource_name="my_endpoint",subscription_id="my_subscription"} 1
# HELP traffic_manager_profile_up Status of the traffic manager profile
# TYPE traffic_manager_profile_up gauge
traffic_manager_profile_up{resource_group="my_rg",resource_name="my_profile",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
