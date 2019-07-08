package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alekssaul/kube-tags2iaas/pkg/appsettings"
	azure "github.com/alekssaul/kube-tags2iaas/pkg/azure"
	"github.com/alekssaul/kube-tags2iaas/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Initializing the application")
	//AppSettings, err := appsettings.GetAppSettings()
	_, err := appsettings.GetAppSettings()
	if err != nil {
		log.Fatalf("Could not initialize the application. %s ", err)
	}

	err = kubernetes.GetUntaggedNodes()
	if err != nil {
		log.Fatalf("Could not get list of Nodes from Kubernetes. %s ", err)
	}

	// err = TagCloudVM("minikube", AppSettings)
	// if err != nil {
	// 	log.Errorf("Could not Tag the Node. %s ", err)
	// }
	os.Exit(0)

}

// TagCloudVM - tags the VM in the cloud
func TagCloudVM(vmName string, settings appsettings.AppSettings) (err error) {
	// Get VM from Azure
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()
	vm, err := azure.GetVM(ctx, vmName)
	if err != nil {
		return fmt.Errorf("Error Retrieving VM object %s from Cloud: %v", vmName, err)
	}

	// Update the tag
	err = azure.UpdateVMTags(ctx, &vm, settings)
	if err != nil {
		return fmt.Errorf("Error Tagging VM object: %v", err)
	}

	return nil
}
