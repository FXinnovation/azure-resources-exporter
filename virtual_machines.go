package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

// VirtualMachinesClient is the client implementation to VirtualMachines API
type VirtualMachinesClient struct {
	Session *AzureSession
	Client  *compute.VirtualMachinesClient
}

// VirtualMachines client interface
type VirtualMachines interface {
	GetVirtualMachines() (*[]compute.VirtualMachine, error)
	GetSubscriptionID() string
}

// NewVirtualMachines returns a new Virtual Machines client
func NewVirtualMachines(subscriptionID string) (VirtualMachines, error) {

	if subscriptionID == "" {
		return nil, errors.New("Invalid subscription ID")
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize authorizer")
	}

	session := AzureSession{
		SubscriptionID: subscriptionID,
		Authorizer:     authorizer,
	}

	vmClient := compute.NewVirtualMachinesClient(session.SubscriptionID)
	vmClient.Authorizer = session.Authorizer

	return &VirtualMachinesClient{
		Session: &session,
		Client:  &vmClient,
	}, nil
}

// GetSubscriptionID return the client's Subscription ID
func (vmClient *VirtualMachinesClient) GetSubscriptionID() string {
	return vmClient.Session.SubscriptionID
}

// GetVirtualMachines fetch Virtual Machines with instance status
func (vmClient *VirtualMachinesClient) GetVirtualMachines() (*[]compute.VirtualMachine, error) {
	var vmList []compute.VirtualMachine

	for it, err := vmClient.Client.ListAllComplete(context.Background(), "false"); it.NotDone(); err = it.Next() {
		if err != nil {
			return nil, err
		}

		tempVM := it.Value()
		labels := ParseResourceLabels(*tempVM.ID)

		vm, err := vmClient.Client.Get(context.Background(), labels["resource_group"], labels["resource_name"], compute.InstanceView)
		if err != nil {
			return nil, err
		}

		vmList = append(vmList, vm)
	}

	return &vmList, nil
}
