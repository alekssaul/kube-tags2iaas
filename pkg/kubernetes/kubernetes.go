package kubernetes

// Initial code in https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
// helped me get started

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetUntaggedNodes - returns nodes that are do not contain the tags we are looking for
func GetUntaggedNodes() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	for {
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		log.Debugf("There are %d nodes in the cluster\n", len(nodes.Items))

		for _, node := range nodes.Items {
			log.Debugf("NodeName: %s", node.Name)
			log.Debugf("Annotations: %v", node.Annotations)
		}

		return nil

		// // Examples for error handling:
		// // - Use helper functions like e.g. errors.IsNotFound()
		// // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		// _, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod not found\n")
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found pod\n")
		// }

		//time.Sleep(10 * time.Second)
	}
}
