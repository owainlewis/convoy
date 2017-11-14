package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/owainlewis/convoy/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var config = flag.String("config", "", "Path to a kubeconfig file")

func main() {
	// send logs to stderr so we can use 'kubectl logs'
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	glog.Info("Running controller")

	clientset, err := buildClient(*config)

	if err != nil {
		glog.Errorf("Failed to build clientset: %s", err)
		return
	}

	controller.NewConvoyController(clientset)

	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})

	if err != nil {
		glog.Errorf("Failed to retrieve pods: %v", err)
		return
	}

	for _, p := range pods.Items {
		glog.V(3).Infof("Found pods: %s/%s", p.Namespace, p.Name)
	}
}

// Build a Kubernetes client.
// Use either local or cluster config depending on conf value
func buildClient(conf string) (*kubernetes.Clientset, error) {
	config, err := getConfig(conf)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return rest.InClusterConfig()
}
