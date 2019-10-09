package network

import (
	"citihub.com/azureutil"
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// Vnets

func getVnetClient() network.VirtualNetworksClient {
	vnetClient := network.NewVirtualNetworksClient(azureutil.GetAzureSubscriptionID())

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		vnetClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return vnetClient
}

// CreateVirtualNetwork creates a virtual network
func CreateVirtualNetwork(ctx context.Context, vnetName string) (vnet network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		network.VirtualNetwork{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
			},
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vnetClient.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(vnetClient)
}

// CreateVirtualNetworkAndSubnets creates a virtual network with two subnets
func CreateVirtualNetworkAndSubnets(ctx context.Context, vnetName, subnet1Name, subnet2Name string, tags map[string]*string) (vnet network.VirtualNetwork, err error) {
	vnetClient := getVnetClient()
	future, err := vnetClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		network.VirtualNetwork{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr(subnet1Name),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.0.0.0/16"),
						},
					},
					{
						Name: to.StringPtr(subnet2Name),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.1.0.0/16"),
						},
					},
				},
			},
			Tags: tags,
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vnetClient.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(vnetClient)
}

// DeleteVirtualNetwork deletes a virtual network given an existing virtual network
func DeleteVirtualNetwork(ctx context.Context, vnetName string) (result network.VirtualNetworksDeleteFuture, err error) {
	vnetClient := getVnetClient()
	return vnetClient.Delete(ctx, azureutil.GetAzureResourceGP(), vnetName)
}
