package client

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	scheme     = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
	restConfig *rest.Config
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
}

func K8sClientOrDie() *kubernetes.Clientset {
	var err error
	restConfig, err = Config("", "")
	if err != nil {
		setupLog.Error(err, "Failed to get config")
		os.Exit(1)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		setupLog.Error(err, "Failed to get k8s clientset")
		os.Exit(1)
	}

	return clientset
}

func GetRestConfig() *rest.Config {
	return restConfig
}

func Config(kubeconfig, masterURL string) (*rest.Config, error) {
	// If a flag is specified with the config location, use that
	if len(kubeconfig) > 0 {
		return clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	// If an env variable is specified with the config locaiton, use that
	if len(os.Getenv("KUBECONFIG")) > 0 {
		fmt.Printf("Using kubeconfig from env var: %s\n", os.Getenv("KUBECONFIG"))
		return clientcmd.BuildConfigFromFlags(masterURL, os.Getenv("KUBECONFIG"))
	}
	// If no explicit location, try the in-cluster config
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	// If no in-cluster config, try the default location in the user's home directory
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags(
			"", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not locate a kubeconfig")
}
