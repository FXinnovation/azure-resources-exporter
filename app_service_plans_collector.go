package main

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// AppServicePlansCollector collect AppServicePlans metrics
type AppServicePlansCollector struct {
	appServicePlans AppServicePlans
}

// NewAppServicePlansCollector returns the collector
func NewAppServicePlansCollector(session *AzureSession) *AppServicePlansCollector {
	appServicePlans := NewAppServicePlans(session)

	return &AppServicePlansCollector{
		appServicePlans: appServicePlans,
	}
}

// Describe to satisfy the collector interface.
func (c *AppServicePlansCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("AppServicePlansCollector", "dummy", nil, nil)
}

// Collect metrics from App Service Plans API
func (c *AppServicePlansCollector) Collect(ch chan<- prometheus.Metric) {

	planList, err := c.appServicePlans.GetAppServicePlans()
	if err != nil {
		log.Errorf("Failed to get app service plan list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectPlanUp(ch, planList)
}

// CollectPlanUp converts AppServicePlans status as a metric
func (c *AppServicePlansCollector) CollectPlanUp(ch chan<- prometheus.Metric, planList *[]web.AppServicePlan) {

	for _, plan := range *planList {
		up := 0.0
		if plan.AppServicePlanProperties.Status == web.StatusOptionsReady {
			up = 1
		}

		labels, err := ParseResourceID(*plan.ID)
		if err != nil {
			log.Errorf("Skipping AppServicePlan: %s", err)
			continue
		}

		labels["subscription_id"] = c.appServicePlans.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("app_service_plan_up", "Status of the app service plan", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}
