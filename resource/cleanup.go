package resource

import (
	"context"
	"log"

	"citihub.com/azureutil"
)

// Cleanup deletes the rescource group created for the sample
func Cleanup(ctx context.Context) error {
	log.Println("deleting resources")
	_, err := DeleteGroup(ctx, azureutil.GetAzureResourceGP())
	return err
}
