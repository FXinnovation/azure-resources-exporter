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
	virtualMachinesClient := compute.NewVirtualMachinesClient(session.SubscriptionID)
	virtualMachinesClient.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &VirtualMachinesClient{
		Session:   session,
		Client:    &virtualMachinesClient,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (client *VirtualMachinesClient) GetSubscriptionID() string {
	return client.Session.SubscriptionID
}

// GetVirtualMachines fetch Virtual Machines with instance status
func (client *VirtualMachinesClient) GetVirtualMachines() (*[]compute.VirtualMachine, error) {
	var vmList []compute.VirtualMachine

	ressources, err := client.Resources.GetResources(virtualMachinesResourceType)
	if err != nil {
		return nil, err
	}

	for _, ressource := range *ressources {
		labels, err := ParseResourceLabels(*ressource.ID)
		if err != nil {
			log.Errorf("Skipping virtual machine: %s", err)
			continue
		}

		vm, err := client.Client.Get(context.Background(), labels["resource_group"], *ressource.Name, compute.InstanceView)
		if err != nil {
			return nil, err
		}

		vmList = append(vmList, vm)
	}

	return &vmList, nil
}
