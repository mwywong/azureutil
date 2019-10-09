package storage

import (
	"citihub.com/azureutil"
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func getStorageAccountsClient() storage.AccountsClient {
	storageAccountsClient := storage.NewAccountsClient(azureutil.GetAzureSubscriptionID())
	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		storageAccountsClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return storageAccountsClient
}

// CreateStorageAccount starts creation of a new storage account and waits for
// the account to be created.
func CreateStorageAccount(ctx context.Context, accountName, accountGroupName string, tags map[string]*string) (storage.Account, error) {
	var s storage.Account
	storageAccountsClient := getStorageAccountsClient()

	result, err := storageAccountsClient.CheckNameAvailability(
		ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(accountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		return s, fmt.Errorf("storage account check-name-availability failed: %v", err)
	}

	if *result.NameAvailable != true {
		return s, fmt.Errorf(
			"storage account name [%s] not available: %v\nserver message: %v",
			accountName, err, *result.Message)
	}

	future, err := storageAccountsClient.Create(
		ctx,
		accountGroupName,
		accountName,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Kind:     storage.Storage,
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
				EnableHTTPSTrafficOnly: to.BoolPtr(true),
			},
			Tags: tags,
		})

	if err != nil {
		return s, fmt.Errorf("failed to start creating storage account: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, storageAccountsClient.Client)
	if err != nil {
		return s, fmt.Errorf("failed to finish creating storage account: %v", err)
	}

	return future.Result(storageAccountsClient)
}

// DeleteStorageAccount deletes an existing storate account
func DeleteStorageAccount(ctx context.Context, accountName, accountGroupName string) (autorest.Response, error) {
	storageAccountsClient := getStorageAccountsClient()
	return storageAccountsClient.Delete(ctx, accountGroupName, accountName)
}

// GetAccountKeys gets the storage account keys
func GetAccountKeys(ctx context.Context, accountName, accountGroupName string) (storage.AccountListKeysResult, error) {
	accountsClient := getStorageAccountsClient()
	return accountsClient.ListKeys(ctx, accountGroupName, accountName)
}

// GetAccountPrimaryKey return the primary key
func GetAccountPrimaryKey(ctx context.Context, accountName, accountGroupName string) string {
	response, err := GetAccountKeys(ctx, accountName, accountGroupName)
	if err != nil {
		log.Fatalf("failed to list keys: %v", err)
	}
	return *(((*response.Keys)[0]).Value)
}
