package main

import (
	"flag"
	"log"
	"time"

	kubeflow "github.com/StatCan/kubeflow-controller/pkg/generated/clientset/versioned"
	informers "github.com/StatCan/kubeflow-controller/pkg/generated/informers/externalversions"
	"github.com/statcan/prob-notebook-controller/pkg/controller"
	"github.com/statcan/prob-notebook-controller/pkg/signals"
	istio "istio.io/client-go/pkg/clientset/versioned"
	istioinformers "istio.io/client-go/pkg/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL      string
	kubeconfig     string
)

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}

func main() {
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("error building kubeconfig: %v", err)
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("error building kubernetes clientset: %v", err)
	}

	istioclient, err := istio.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("error building istio client: %v", err)
	}

	kubeflowclient, err := kubeflow.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("error building kubeflow client: %v", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, time.Second*30)
	istioInformerFactory := istioinformers.NewSharedInformerFactory(istioclient, time.Second*30)
	kubeflowInformerFactory := informers.NewSharedInformerFactory(kubeflowclient, time.Second*30)

	ctlr := controller.NewController(
		kubeclient,
		istioclient,
		kubeflowclient,
		kubeflowInformerFactory.Kubeflow().V1().Notebooks(),
		istioInformerFactory.Security().V1beta1().AuthorizationPolicies(),
	)

	kubeInformerFactory.Start(stopCh)
	istioInformerFactory.Start(stopCh)
	kubeflowInformerFactory.Start(stopCh)

	if err = ctlr.Run(2, stopCh); err != nil {
		log.Fatalf("error running controller: %v", err)
	}
}
