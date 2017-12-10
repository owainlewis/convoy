package main

import (
	"flag"
	"os"
	"time"

	glog "github.com/golang/glog"
	slack "github.com/nlopes/slack"
	controller "github.com/owainlewis/convoy/pkg/controller"
	notifier "github.com/owainlewis/convoy/pkg/notifier"
	informers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

var config = flag.String("config", "", "Path to a kubeconfig file")

func main() {

	flag.Parse()

	glog.Info("Running controller")

	client, err := buildClient(*config)

	if err != nil {
		glog.Errorf("Failed to build clientset: %s", err)
		return
	}

	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)

	slackClient := slack.New(os.Getenv("SLACK_TOKEN"))
	notifier := notifier.NewSlackNotifier(slackClient, "convoyk8s")

	ctrl := controller.NewConvoyController(
		client,
		sharedInformers.Core().V1().Events(),
		notifier)

	stopCh := make(chan struct{})

	defer close(stopCh)

	sharedInformers.Start(stopCh)
	ctrl.Run(stopCh)
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
