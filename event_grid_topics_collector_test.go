package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2019-02-01-preview/eventgrid"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedEventGridTopics struct {
	mock.Mock
}

func (mock *MockedEventGridTopics) GetEventGridTopics() (*[]eventgrid.Topic, error) {
	args := mock.Called()
	return args.Get(0).(*[]eventgrid.Topic), args.Error(1)
}

func (mock *MockedEventGridTopics) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func TestNewEventGridTopicsCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewEventGridTopicsCollector(session)
}

func TestCollectEventGridTopic_Error(t *testing.T) {
	eventGridTopics := MockedEventGridTopics{}
	collector := EventGridTopicsCollector{
		eventGridTopics: &eventGridTopics,
	}
	prometheus.MustRegister(&collector)

	var egtList []eventgrid.Topic
	eventGridTopics.On("GetEventGridTopics").Return(&egtList, errors.New("Unit test Error"))
	eventGridTopics.On("GetSubscriptionID").Return("my_subscription")

	req := httptest.NewRequest("GET", "/webhook", nil)
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollectEventGridTopic_Up(t *testing.T) {
	eventGridTopics := MockedEventGridTopics{}
	collector := EventGridTopicsCollector{
		eventGridTopics: &eventGridTopics,
	}

	var egtList []eventgrid.Topic
	id := "/subscriptions/my_subscription/resourceGroups/my_rg/providers/Microsoft.EventGrid/topics/my_topic"
	tagValue := "Value"
	resourceType := "Microsoft.EventGrid/topics"

	egtList = append(egtList, eventgrid.Topic{
		TopicProperties: &eventgrid.TopicProperties{
			ProvisioningState: eventgrid.TopicProvisioningStateSucceeded,
		},
		ID: &id,
		Tags: map[string]*string{
			"Key": &tagValue,
		},
		Type: &resourceType,
	})

	eventGridTopics.On("GetEventGridTopics").Return(&egtList, nil)
	eventGridTopics.On("GetSubscriptionID").Return("my_subscription")

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
azure_tag_info{resource_group="my_rg",resource_name="my_topic",resource_type="Microsoft.EventGrid/topics",subscription_id="my_subscription",tag_key="Value"} 1
# HELP event_grid_topic_up Provisionning state of the event grid topic
# TYPE event_grid_topic_up gauge
event_grid_topic_up{resource_group="my_rg",resource_name="my_topic",subscription_id="my_subscription"} 1
`
	if rr.Body.String() != want {
		t.Errorf("Unexpected body: got %v, want %v", rr.Body.String(), want)
	}
}
