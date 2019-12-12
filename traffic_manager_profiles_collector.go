package main

import (
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// TrafficManagerProfilesCollector collect Traffic Manager Profiles metrics
type TrafficManagerProfilesCollector struct {
	trafficManagerProfiles TrafficManagerProfiles
}

// NewTrafficManagerProfilesCollector returns the collector
func NewTrafficManagerProfilesCollector(session *AzureSession) *TrafficManagerProfilesCollector {
	trafficManagerProfiles := NewTrafficManagerProfiles(session)

	return &TrafficManagerProfilesCollector{
		trafficManagerProfiles: trafficManagerProfiles,
	}
}

// Describe to satisfy the collector interface.
func (c *TrafficManagerProfilesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("TrafficManagerProfilesCollector", "dummy", nil, nil)
}

// Collect metrics from TrafficManagerProfiles API
func (c *TrafficManagerProfilesCollector) Collect(ch chan<- prometheus.Metric) {

	profileList, err := c.trafficManagerProfiles.GetTrafficManagerProfiles()
	if err != nil {
		log.Errorf("Failed to get traffic manager profile list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectTrafficManagerProfileUp(ch, profileList)
}

// CollectTrafficManagerProfileUp converts Traffic Manager Profiles and endpoint status as a metric
func (c *TrafficManagerProfilesCollector) CollectTrafficManagerProfileUp(ch chan<- prometheus.Metric, profileList *[]trafficmanager.Profile) {

	for _, profile := range *profileList {
		up := 0.0
		if profile.ProfileProperties.ProfileStatus == trafficmanager.ProfileStatusEnabled {
			up = 1
		}

		labels, err := ParseResourceID(*profile.ID)
		if err != nil {
			log.Errorf("Skipping traffic manager profile: %s", err)
			continue
		}

		labels["subscription_id"] = c.trafficManagerProfiles.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("traffic_manager_profile_up", "Status of the traffic manager profile", nil, labels),
			prometheus.GaugeValue,
			up,
		)

		ExportAzureTagInfo(ch, profile.Tags, profile.Type, labels)

		for _, endpoint := range *profile.Endpoints {
			endpointUp := 0.0
			if endpoint.EndpointProperties.EndpointStatus == trafficmanager.EndpointStatusEnabled {
				endpointUp = 1
			}

			endpointLabels, err := ParseResourceID(*endpoint.ID)
			if err != nil {
				log.Errorf("Skipping traffic manager profile endpoint: %s", err)
				continue
			}

			endpointLabels["subscription_id"] = c.trafficManagerProfiles.GetSubscriptionID()

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("traffic_manager_profile_endpoint_up", "Status of the traffic manager profile endpoint", nil, endpointLabels),
				prometheus.GaugeValue,
				endpointUp,
			)

			ExportAzureTagInfo(ch, profile.Tags, endpoint.Type, endpointLabels)
		}
	}
}
