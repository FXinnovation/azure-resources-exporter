package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedWebApps struct {
	mock.Mock
}

func (mock *MockedWebApps) GetWebApps() (*[]web.Site, error) {
	args := mock.Called()
	return args.Get(0).(*[]web.Site), args.Error(1)
}

func (mock *MockedWebApps) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewWebAppsCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewWebAppsCollector(session)
}

func TestCollectApp_Error(t *testing.T) {
	webApps := MockedWebApps{}
	collector := WebAppsCollector{
		webApps: &webApps,
	}
	prometheus.MustRegister(&collector)

	var appList []web.Site
	webApps.On("GetWebApps").Return(&appList, errors.New("Unit test Error"))
	webApps.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectApp_Up(t *testing.T) {
	webApps := MockedWebApps{}
	collector := WebAppsCollector{
		webApps: &webApps,
	}

	var appList []web.Site
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Web/site/my_app"
	state := "Running"

	appList = append(appList, web.Site{
		SiteProperties: &web.SiteProperties{
			State: &state,
		},
		ID: &id,
	})

	webApps.On("GetWebApps").Return(&appList, nil)
	webApps.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP web_app_up Status of the web app
# TYPE web_app_up gauge
web_app_up{resource_group="my_rg",resource_name="my_app",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
