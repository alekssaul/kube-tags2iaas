package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alekssaul/kube-tags2iaas/pkg/appsettings"
	azure "github.com/alekssaul/kube-tags2iaas/pkg/azure"
	"github.com/alekssaul/kube-tags2iaas/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Println("Initializing the application")
	AppSettings, err := appsettings.GetAppSettings()
	if err != nil {
		log.Fatalf("Could not initialize the application. %s ", err)
	}

	nodes2tag, err := kubernetes.GetUntaggedNodes(AppSettings)
	if err != nil {
		log.Fatalf("Could not get list of Nodes from Kubernetes. %s ", err)
	}
	log.Debugf("nodes2tag: %v", nodes2tag)

	for _, node := range nodes2tag {
		log.Printf("Node %s needs to be tagged, calling the Cloud API", node)
		err = TagCloudVM(node, AppSettings)
		if err != nil {
			log.Errorf("Could not Tag the Node %s. %s ", node, err)
		}

		err = kubernetes.SetTagonNode(node, AppSettings)
		if err != nil {
			log.Errorf("Could not annotate the Node %s on K8s API . %s ", node, err)
		}
	}

	time.Sleep(60 * time.Minute)

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
