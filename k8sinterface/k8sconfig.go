package k8sinterface

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

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
var clientConfigAPI *clientcmdapi.Config

// KubernetesApi -
type KubernetesApi struct {
	ApiExtensionsClient clientset.Interface
	KubernetesClient    kubernetes.Interface
	DynamicClient       dynamic.Interface
	DiscoveryClient     discovery.DiscoveryInterface
	Context             context.Context
	K8SConfig           *restclient.Config
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

	apiExtensionsClient, err := clientset.NewForConfig(k8sConfig)
	if err != nil {
		logger.L().Fatal("failed to initialize a new discovery client", helpers.Error(err))
	}

	restclient.SetDefaultWarningHandler(restclient.NoWarnings{})
	InitializeMapResources(discoveryClient)

	return &KubernetesApi{
		ApiExtensionsClient: apiExtensionsClient,
		KubernetesClient:    kubernetesClient,
		DynamicClient:       dynamicClient,
		DiscoveryClient:     discoveryClient,
		Context:             context.Background(),
		K8SConfig:           k8sConfig,
	}
}
func (k8sAPI *KubernetesApi) GetKubernetesClient() kubernetes.Interface {
	return k8sAPI.KubernetesClient
}
func (k8sAPI *KubernetesApi) GetDynamicClient() dynamic.Interface {
	return k8sAPI.DynamicClient
}
func (k8sAPI *KubernetesApi) GetDiscoveryClient() discovery.DiscoveryInterface {
	return k8sAPI.DiscoveryClient
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

func GetContext() *clientcmdapi.Context {
	kubeConfig := GetConfig()
	if kubeConfig == nil {
		return nil
	}

	contextName := clusterContextName
	if contextName == "" {
		// if context name is not set, use the current context
		contextName = kubeConfig.CurrentContext
	}

	if context, exist := kubeConfig.Contexts[contextName]; exist && context != nil {
		// return the context
		return context
	}
	return nil
}

func SetClientConfigAPI(conf *clientcmdapi.Config) {
	clientConfigAPI = conf
}

func IsConnectedToCluster() bool {
	if K8SConfig == nil {
		if err := LoadK8sConfig(); err != nil {
			SetConnectedToCluster(false)
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
	if clientConfigAPI != nil {
		return clientConfigAPI
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{CurrentContext: clusterContextName})
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil
	}

	// set the config to the global variable
	SetClientConfigAPI(&config)

	return &config
}

// GetDefaultNamespace returns the default namespace for the current context
func GetDefaultNamespace() string {

	if context := GetContext(); context != nil {
		return context.Namespace
	}

	// return default namespace in case the context is not available
	return "default"
}

// GetCluster returns a pointer to the clientcmdapi Cluster object of the current context
func GetCluster() *clientcmdapi.Cluster {
	config := GetConfig()
	if config == nil {
		return nil
	}

	if context, exist := config.Clusters[GetContextName()]; exist && context != nil {
		// return the cluster as based on the context
		return context
	}
	return nil

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

// GetK8sConfigClusterServerName get the server name of desired cluster context
func GetK8sConfigClusterServerName() string {

	config := GetConfig()
	if config == nil {
		return ""
	}

	if context, exist := config.Clusters[GetContextName()]; exist && context != nil {
		// return the server name of the context
		return context.Server
	}

	// return current context in case the server name is not available
	return ConfigClusterServerName
}

func SetConnectedToCluster(isConnected bool) {
	connectedToCluster = isConnected
}
