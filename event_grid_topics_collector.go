package main

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2019-02-01-preview/eventgrid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// EventGridTopicsCollector collect EventGridTopics metrics
type EventGridTopicsCollector struct {
	eventGridTopics EventGridTopics
}

// NewEventGridTopicsCollector returns the collector
func NewEventGridTopicsCollector(session *AzureSession) *EventGridTopicsCollector {
	eventGridTopics := NewEventGridTopics(session)

	return &EventGridTopicsCollector{
		eventGridTopics: eventGridTopics,
	}
}

// Describe to satisfy the collector interface.
func (c *EventGridTopicsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("EventGridTopicsCollector", "dummy", nil, nil)
}

// Collect metrics from EventGridTopics API
func (c *EventGridTopicsCollector) Collect(ch chan<- prometheus.Metric) {

	egtList, err := c.eventGridTopics.GetEventGridTopics()
	if err != nil {
		log.Errorf("Failed to get Event Grid Topic list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectEventGridTopicUp(ch, egtList)
}

// CollectEventGridTopicUp converts EventGridTopics state as a metric
func (c *EventGridTopicsCollector) CollectEventGridTopicUp(ch chan<- prometheus.Metric, egtList *[]eventgrid.Topic) {

	for _, egt := range *egtList {
		up := 0.0
		if egt.TopicProperties.ProvisioningState == eventgrid.TopicProvisioningStateSucceeded {
			up = 1
		}

		labels, err := ParseResourceID(*egt.ID)
		if err != nil {
			log.Errorf("Skipping Event Grid Topic: %s", err)
			continue
		}

		labels["subscription_id"] = c.eventGridTopics.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("event_grid_topic_up", "Provisionning state of the event grid topic", nil, labels),
			prometheus.GaugeValue,
			up,
		)

		ExportAzureTagInfo(ch, egt.Tags, egt.Type, labels)
	}
}
