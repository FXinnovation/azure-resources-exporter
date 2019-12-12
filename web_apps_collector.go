package main

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// AppUpState is the state of an App that is UP
const AppUpState = "Running"

// WebAppsCollector collect WebApps metrics
type WebAppsCollector struct {
	webApps WebApps
}

// NewWebAppsCollector returns the collector
func NewWebAppsCollector(session *AzureSession) *WebAppsCollector {
	webApps := NewWebApps(session)

	return &WebAppsCollector{
		webApps: webApps,
	}
}

// Describe to satisfy the collector interface.
func (c *WebAppsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("WebAppsCollector", "dummy", nil, nil)
}

// Collect metrics from Web Apps API
func (c *WebAppsCollector) Collect(ch chan<- prometheus.Metric) {

	appList, err := c.webApps.GetWebApps()
	if err != nil {
		log.Errorf("Failed to get web app list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectAppUp(ch, appList)
}

// CollectAppUp converts WebApps state as a metric
func (c *WebAppsCollector) CollectAppUp(ch chan<- prometheus.Metric, appList *[]web.Site) {

	for _, app := range *appList {
		up := 0.0
		if *app.SiteProperties.State == AppUpState {
			up = 1
		}

		labels, err := ParseResourceID(*app.ID)
		if err != nil {
			log.Errorf("Skipping WebApp: %s", err)
			continue
		}

		labels["subscription_id"] = c.webApps.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("web_app_up", "Status of the web app", nil, labels),
			prometheus.GaugeValue,
			up,
		)

		ExportAzureTagInfo(ch, app.Tags, app.Type, labels)
	}
}
