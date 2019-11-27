package main

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// VirtualMachinesCollector collect Virtual Machines metrics
type VirtualMachinesCollector struct {
	virtualMachines VirtualMachines
}

// NewVirtualMachinesCollector returns the collector
func NewVirtualMachinesCollector(subscriptionID string) (*VirtualMachinesCollector, error) {
	virtualMachines, err := NewVirtualMachines(subscriptionID)

	if err != nil {
		return nil, err
	}

	return &VirtualMachinesCollector{
		virtualMachines: virtualMachines,
	}, nil
}

// Describe to satisfy the collector interface.
func (c *VirtualMachinesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect metrics from Virtual Machine API
func (c *VirtualMachinesCollector) Collect(ch chan<- prometheus.Metric) {

	vmList, err := c.virtualMachines.GetVirtualMachines()

	if err != nil {
		log.Errorf("Failed to get virtual machines list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectInstanceUp(ch, vmList)
}

// CollectInstanceUp converts VM instance running status as a metric
func (c *VirtualMachinesCollector) CollectInstanceUp(ch chan<- prometheus.Metric, vmList *[]compute.VirtualMachine) {

	for _, vm := range *vmList {
		up := 0.0
		for _, status := range *vm.VirtualMachineProperties.InstanceView.Statuses {
			if *status.Code == "PowerState/running" {
				up = 1
				break
			}
		}

		labels, err := ParseResourceLabels(*vm.ID)

		if err != nil {
			log.Errorf("Skipping virtual machine: %s", err)
			continue
		}

		labels["subscription_id"] = c.virtualMachines.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("virtual_machine_instance_up", "Running status of the virtual machine instance", nil, labels),
			prometheus.GaugeValue,
			up,
		)
	}
}