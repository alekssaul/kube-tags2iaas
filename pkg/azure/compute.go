package azure

// getVMClient and GetVM are referances from ; https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/compute/vm.go

import (
	"context"

	"github.com/alekssaul/kube-tags2iaas/pkg/appsettings"
	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/alekssaul/kube-tags2iaas/pkg/azure/internal/config"
	"github.com/alekssaul/kube-tags2iaas/pkg/azure/internal/iam"
)

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	vmClient.Authorizer = a
	vmClient.AddToUserAgent(config.UserAgent())
	return vmClient
}

// GetVM gets the specified VM info
func GetVM(ctx context.Context, vmName string) (compute.VirtualMachine, error) {
	err := config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err)
	}
	vmClient := getVMClient()
	return vmClient.Get(ctx, config.GroupName(), vmName, compute.InstanceView)
}

// UpdateVMTags updates the VM with tags
func UpdateVMTags(ctx context.Context, vm *compute.VirtualMachine, appSettings appsettings.AppSettings) error {
	// Append Infrastructure tags to existing VM.Tags
	newTags := make(map[string]*string)
	if vm.Tags != nil {
		newTags = vm.Tags
	}

	for i := range appSettings.InfrastructureTags {
		newTags[appSettings.InfrastructureTags[i].Key] = &appSettings.InfrastructureTags[i].Value
	}

	vm.Tags = newTags

	vmClient := getVMClient()
	future, err := vmClient.CreateOrUpdate(ctx, config.GroupName(), *vm.Name, *vm)
	if err != nil {
		return fmt.Errorf("cannot update vm: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get the vm create or update future response: %v", err)
	}

	log.Printf("Updated VM: %s with tags %s", *vm.Name, vm.Tags)

	return nil

}
