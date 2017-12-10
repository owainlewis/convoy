package main

import (
	"flag"
	"os"
	"time"

	glog "github.com/golang/glog"
	config "github.com/owainlewis/convoy/pkg/config"
	controller "github.com/owainlewis/convoy/pkg/controller"
	dispatch "github.com/owainlewis/convoy/pkg/dispatch"
	informers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

var conf = flag.String("config", "", "Path to config YAML")
var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {
	flag.Parse()

	glog.Info("Starting controller")

	client, err := buildClient(*kubeconfig)
	if err != nil {
		glog.Errorf("Failed to build clientset: %s", err)
		return
	}

	c, err := getConfig(*conf)
	if err != nil {
		glog.Fatalf("Failed to load configuration %s", err)
	}

	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)

	dispatcher := getDispatcher(*c)

	ctrl := controller.NewConvoyController(
		client,
		sharedInformers.Core().V1().Events(),
		dispatcher,
		*c)

	stopCh := make(chan struct{})

	defer close(stopCh)

	sharedInformers.Start(stopCh)
	ctrl.Run(stopCh)
}

// Build a Kubernetes client.
// Use either local or cluster config depending on conf value
func buildClient(conf string) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(conf)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getDispatcher(config config.Config) dispatch.Notifier {
	if config.Slack.Enabled {
		slackToken := os.Getenv("SLACK_TOKEN")
		slackChannel := config.Slack.Channel
		return dispatch.NewSlackNotifier(slackToken, slackChannel)
	}
	return dispatch.NewConsoleNotifier()
}

func getKubeConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return rest.InClusterConfig()
}

func getConfig(conf string) (*config.Config, error) {
	defaultConfigPath := "config.yml"
	if conf == "" {
		return config.FromFile(defaultConfigPath)
	}
	return config.FromFile(conf)
}
