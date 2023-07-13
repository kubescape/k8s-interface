package k8sinterface

import (
	"context"
	"fmt"
	"strings"

	logger "github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// DO NOT REMOVE - load cloud providers auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var connectedToCluster = true
var clusterContextName = ""
var ConfigClusterServerName = ""
var K8SGitServerVersion = ""

// K8SConfig pointer to k8s config
var K8SConfig *restclient.Config

// KubernetesApi -
type KubernetesApi struct {
	KubernetesClient kubernetes.Interface
	DynamicClient    dynamic.Interface
	DiscoveryClient  discovery.DiscoveryInterface
	Context          context.Context
	K8SConfig        *restclient.Config
}

// NewKubernetesApi -
func NewKubernetesApi() *KubernetesApi {
	var kubernetesClient *kubernetes.Clientset
	var err error

	if !IsConnectedToCluster() {
		logger.L().Fatal("failed to load kubernetes config: no configuration has been provided, try setting KUBECONFIG environment variable")
	}

	k8sConfig := GetK8sConfig()

	kubernetesClient, err = kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new kubernetes client", helpers.Error(err))
	}

	dynamicClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new dynamic client", helpers.Error(err))
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new discovery client", helpers.Error(err))
	}
	restclient.SetDefaultWarningHandler(restclient.NoWarnings{})
	InitializeMapResources(discoveryClient)

	return &KubernetesApi{
		KubernetesClient: kubernetesClient,
		DynamicClient:    dynamicClient,
		DiscoveryClient:  discoveryClient,
		Context:          context.Background(),
		K8SConfig:        k8sConfig,
	}
}

// RunningIncluster whether running in cluster
var RunningIncluster bool

// LoadK8sConfig load config from local file or from cluster
func LoadK8sConfig() error {
	kubeconfig, err := config.GetConfigWithContext(clusterContextName)
	if err != nil {
		return fmt.Errorf("failed to load kubernetes config: %s", strings.ReplaceAll(err.Error(), "KUBERNETES_MASTER", "KUBECONFIG"))
	}
	if _, err := restclient.InClusterConfig(); err == nil {
		RunningIncluster = true
	} else {
		RunningIncluster = false
	}

	K8SConfig = kubeconfig
	return nil
}

// GetK8sConfig get config. load if not loaded yet
func GetK8sConfig() *restclient.Config {
	if !IsConnectedToCluster() {
		return nil
	}
	return K8SConfig
}

// DEPRECATED
func GetCurrentContext() *clientcmdapi.Context {
	if kubeConfig := GetConfig(); kubeConfig != nil {
		if clusterContextName != "" {
			if c, ok := kubeConfig.Contexts[clusterContextName]; ok {
				return c
			}
		}
		// if context name is not set, return the current context
		return kubeConfig.Contexts[kubeConfig.CurrentContext]
	}
	return nil
}

func IsConnectedToCluster() bool {
	if K8SConfig == nil {
		if err := LoadK8sConfig(); err != nil {
			connectedToCluster = false
		}
	}
	return connectedToCluster
}

func GetContextName() string {
	if clusterContextName != "" {
		return clusterContextName
	}

	if config := GetConfig(); config != nil {
		return config.CurrentContext
	}

	return ""
}

// get config from ~/.kube/config
func GetConfig() *clientcmdapi.Config {

	if !connectedToCluster {
		return nil
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{CurrentContext: clusterContextName})
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil
	}
	return &config
}

// GetDefaultNamespace returns the default namespace for the current context
func GetDefaultNamespace() string {
	defaultNamespace := "default"
	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return defaultNamespace
	}

	tempClusterContextName := clusterContextName
	if tempClusterContextName == "" {
		tempClusterContextName = clientCfg.CurrentContext
	}
	apiContext, ok := clientCfg.Contexts[tempClusterContextName]
	if !ok || apiContext == nil {
		return defaultNamespace
	}
	namespace := apiContext.Namespace
	if apiContext.Namespace == "" {
		namespace = defaultNamespace
	}
	return namespace
}

// SetClusterContextName set the name of desired cluster context. The package will use this name when loading the context
func SetClusterContextName(contextName string) {
	clusterContextName = contextName
}

func SetK8SGitServerVersion(K8SGitServerVersionInput string) {
	K8SGitServerVersion = K8SGitServerVersionInput
}

func SetConfigClusterServerName(contextName string) {
	ConfigClusterServerName = contextName
}

func GetK8sConfigClusterServerName() string {
	if ConfigClusterServerName == "" {
		config := GetConfig()
		if _, exist := config.Clusters[config.CurrentContext]; !exist {
			ConfigClusterServerName = config.Clusters[config.CurrentContext].Server
		}
	}
	return ConfigClusterServerName
}
