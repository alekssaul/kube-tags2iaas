package kubernetes

// Initial code in https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
// helped me get started

import (
	"github.com/alekssaul/kube-tags2iaas/pkg/appsettings"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetUntaggedNodes - returns nodes that are do not contain the tags we are looking for
func GetUntaggedNodes(appSettings appsettings.AppSettings) (nodes2tag []string, err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nodes2tag, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nodes2tag, err
	}
	for {
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return nodes2tag, err
		}
		log.Debugf("There are %d nodes in the cluster\n", len(nodes.Items))

		for _, node := range nodes.Items {
			log.Debugf("NodeName: %s", node.Name)
			log.Debugf("Annotations: %v", node.Annotations)

			// not very pretty, this will itirate over annotations of all nodes for all annotations
			for i := range appSettings.InfrastructureTags {
				if node.Annotations[appSettings.InfrastructureTags[i].Key] == "" {
					log.Debugf("Annotation: %v, does not exist on node %s", appSettings.InfrastructureTags[i].Key, node.Name)
					log.Debugf("Adding Node %s, to list of nodes to tag", node.Name)
					nodes2tag = append(nodes2tag, node.Name)
					break
				} else if node.Annotations[appSettings.InfrastructureTags[i].Key] != appSettings.InfrastructureTags[i].Value {
					log.Debugf("Annotation: %v, value is %s however we expected %s",
						appSettings.InfrastructureTags[i].Key,
						node.Annotations[appSettings.InfrastructureTags[i].Key],
						appSettings.InfrastructureTags[i].Value,
					)
					log.Debugf("Adding Node %s, to list of nodes to tag", node.Name)
					nodes2tag = append(nodes2tag, node.Name)
					break
				}
			}
		}

		return nodes2tag, nil
	}
}
