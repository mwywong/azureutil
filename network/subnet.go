package network

import (
	"context"
	"fmt"
	"log"

	"citihub.com/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-11-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// VNet Subnets

func getSubnetsClient() network.SubnetsClient {
	subnetsClient := network.NewSubnetsClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		subnetsClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return subnetsClient
}

// CreateVirtualNetworkSubnet creates a subnet in an existing vnet
func CreateVirtualNetworkSubnet(ctx context.Context, vnetName, subnetName string) (subnet network.Subnet, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		subnetName,
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/16"),
			},
		})
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// CreateSubnetWithNetworkSecurityGroup create a subnet referencing a network security group
func CreateSubnetWithNetworkSecurityGroup(ctx context.Context, vnetName, subnetName, addressPrefix, nsgName string) (subnet network.Subnet, err error) {
	nsg, err := GetNetworkSecurityGroup(ctx, nsgName)
	if err != nil {
		return subnet, fmt.Errorf("cannot get nsg: %v", err)
	}

	subnetsClient := getSubnetsClient()
	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		subnetName,
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix:        to.StringPtr(addressPrefix),
				NetworkSecurityGroup: &nsg,
			},
		})
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// DeleteVirtualNetworkSubnet deletes a subnet
func DeleteVirtualNetworkSubnet() {}

// GetVirtualNetworkSubnet returns an existing subnet from a virtual network
func GetVirtualNetworkSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	subnetsClient := getSubnetsClient()
	return subnetsClient.Get(ctx, azureutil.GetAzureResourceGP(), vnetName, subnetName, "")
}
