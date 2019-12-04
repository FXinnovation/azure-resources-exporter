package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/prometheus/common/log"
)

const virtualMachinesResourceType = "Microsoft.Compute/virtualMachines"

// VirtualMachinesClient is the client implementation to VirtualMachines API
type VirtualMachinesClient struct {
	Session   *AzureSession
	Client    *compute.VirtualMachinesClient
	Resources Resources
}

// VirtualMachines client interface
type VirtualMachines interface {
	GetVirtualMachines() (*[]compute.VirtualMachine, error)
	GetSubscriptionID() string
}

// NewVirtualMachines returns a new Virtual Machines client
func NewVirtualMachines(session *AzureSession) VirtualMachines {
	client := compute.NewVirtualMachinesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &VirtualMachinesClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (vc *VirtualMachinesClient) GetSubscriptionID() string {
	return vc.Session.SubscriptionID
}

// GetVirtualMachines fetch Virtual Machines with instance status
func (vc *VirtualMachinesClient) GetVirtualMachines() (*[]compute.VirtualMachine, error) {
	var vmList []compute.VirtualMachine

	resources, err := vc.Resources.GetResources(virtualMachinesResourceType)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		labels, err := ParseResourceLabels(*resource.ID)
		if err != nil {
			log.Errorf("Skipping virtual machine: %s", err)
			continue
		}

		vm, err := vc.Client.Get(context.Background(), labels["resource_group"], *resource.Name, compute.InstanceView)
		if err != nil {
			return nil, err
		}

		vmList = append(vmList, vm)
	}

	return &vmList, nil
}
