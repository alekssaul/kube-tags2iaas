package main

import (
	"context"
	"fmt"

	"time"

	"github.com/alekssaul/kube-tags2iaas/pkg/appsettings"
	azure "github.com/alekssaul/kube-tags2iaas/pkg/azure"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Initializing the application")
	AppSettings, err := appsettings.GetAppSettings()
	if err != nil {
		log.Fatalf("Could not initialize the application. %s ", err)
	}

	err = TagCloudVM("test", AppSettings)
	if err != nil {
		log.Errorf("Could not initialize the application. %s ", err)
	}

}

// TagCloudVM - tags the VM in the cloud
func TagCloudVM(vmName string, settings appsettings.AppSettings) (err error) {
	// Get VM from Azure
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()
	vm, err := azure.GetVM(ctx, "test")
	if err != nil {
		return fmt.Errorf("Error Retrieving VM object from Cloud: %v", err)
	}

	// Update the tag
	err = azure.UpdateVMTags(ctx, &vm, settings)
	if err != nil {
		return fmt.Errorf("Error Tagging VM object: %v", err)
	}

	return nil
}
