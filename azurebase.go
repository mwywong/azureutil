package azureutil

import (
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const azureLocation string = "AZURE_LOCATION"
const azureSubscriptionID string = "AZURE_SUBSCRIPTION_ID"

var testprefix string
var azureResourceGp string

// GetTestPrefix return a random test prefix with test + 6 random characters
func GetTestPrefix() string {
	if testprefix == "" {
		testprefix = "test" + RandStringBytesMaskImprSrcUnsafe(6) + ""
	}
	return testprefix
}

func getFromEnvVar(varName string) string {
	result := os.Getenv(varName)
	if result == "" {
		log.Fatalf("Environment variables \"%v\" is not defined", varName)
	}
	return result

}

//GetAzureResourceGP - Default resourece GP
func GetAzureResourceGP() string {
	if azureResourceGp == "" {
		azureResourceGp = GetTestPrefix() + "resourecGP"
	}
	return azureResourceGp
}

//GetAzureLocation - Default location
func GetAzureLocation() string {
	return getFromEnvVar(azureLocation)
}

//GetAzureSubscriptionID - Return subscriptionID
func GetAzureSubscriptionID() string {
	return getFromEnvVar(azureSubscriptionID)
}

//GetAzureAuthorizer - return an Azure Authorizer
func GetAzureAuthorizer() autorest.Authorizer {
	// create an authorizer from env vars or Azure Managed Service Idenity

	Authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err != nil {
		log.Panicf("Unable to load Azure credential due to %v", err)
	}
	return Authorizer
}
