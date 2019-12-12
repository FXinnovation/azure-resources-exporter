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

type MockedAppServicePlans struct {
	mock.Mock
}

func (mock *MockedAppServicePlans) GetAppServicePlans() (*[]web.AppServicePlan, error) {
	args := mock.Called()
	return args.Get(0).(*[]web.AppServicePlan), args.Error(1)
}

func (mock *MockedAppServicePlans) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewAppServicePlansCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewAppServicePlansCollector(session)
}

func TestCollectPlan_Error(t *testing.T) {
	appServicePlans := MockedAppServicePlans{}
	collector := AppServicePlansCollector{
		appServicePlans: &appServicePlans,
	}
	prometheus.MustRegister(&collector)

	var planList []web.AppServicePlan
	appServicePlans.On("GetAppServicePlans").Return(&planList, errors.New("Unit test Error"))
	appServicePlans.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectPlan_Up(t *testing.T) {
	appServicePlans := MockedAppServicePlans{}
	collector := AppServicePlansCollector{
		appServicePlans: &appServicePlans,
	}

	var planList []web.AppServicePlan
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.Web/serverfarms/my_plan"
	tagValue := "Value"
	resourceType := "Microsoft.Web/serverfarms"

	planList = append(planList, web.AppServicePlan{
		AppServicePlanProperties: &web.AppServicePlanProperties{
			Status: web.StatusOptionsReady,
		},
		ID: &id,
		Tags: map[string]*string{
			"Key": &tagValue,
		},
		Type: &resourceType,
	})

	appServicePlans.On("GetAppServicePlans").Return(&planList, nil)
	appServicePlans.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusOK)
	}

	want := `# HELP app_service_plan_up Status of the app service plan
# TYPE app_service_plan_up gauge
app_service_plan_up{resource_group="my_rg",resource_name="my_plan",subscription_id="my_subscription"} 1
# HELP azure_tag_info Tags of the Azure resource
# TYPE azure_tag_info gauge
azure_tag_info{resource_group="my_rg",resource_name="my_plan",resource_type="Microsoft.Web/serverfarms",subscription_id="my_subscription",tag_key="Value"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
