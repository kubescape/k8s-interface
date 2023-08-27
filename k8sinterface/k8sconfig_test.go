package k8sinterface

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/stretchr/testify/assert"
)

func TestGetK8sConfigClusterServerName(t *testing.T) {
	defer tearDown()

	tests := []struct {
		name           string
		config         *clientcmdapi.Config
		expectedServer string
	}{
		{
			name:           "Config is empty, return default server name",
			config:         &clientcmdapi.Config{},
			expectedServer: ConfigClusterServerName,
		},
		{
			name: "Context name is not set, return current context server name",
			config: &clientcmdapi.Config{
				CurrentContext: "test-context",
				Clusters: map[string]*clientcmdapi.Cluster{
					"test-context": {
						Server: "https://test-server.com",
					},
				},
			},
			expectedServer: "https://test-server.com",
		},
		{
			name: "Context name is set, return context server name",
			config: &clientcmdapi.Config{
				CurrentContext: "test-context",
				Clusters: map[string]*clientcmdapi.Cluster{
					"test-context": {
						Server: "https://test-server.com",
					},
				},
			},
			expectedServer: "https://test-server.com",
		},
		{
			name: "Context name is set, but server name is not available, return default server name",
			config: &clientcmdapi.Config{
				CurrentContext: "test-context",
				Clusters: map[string]*clientcmdapi.Cluster{
					"test-context": {},
				},
			},
			expectedServer: ConfigClusterServerName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetClientConfigAPI(tt.config)
			actualServer := GetK8sConfigClusterServerName()
			assert.Equal(t, tt.expectedServer, actualServer)
		})
	}
}

func TestGetContext(t *testing.T) {
	defer tearDown()
	tests := []struct {
		name           string
		clusterContext string
		config         *clientcmdapi.Config
		expectedCtx    *clientcmdapi.Context
	}{
		{
			name:        "Config is nil, return nil context",
			config:      &clientcmdapi.Config{},
			expectedCtx: nil,
		},
		{
			name: "Cluster context name is not set, return current context",
			config: &clientcmdapi.Config{
				CurrentContext: "test-context",
				Contexts: map[string]*clientcmdapi.Context{
					"test-context": {
						Namespace: "test-namespace",
					},
				},
			},
			expectedCtx: &clientcmdapi.Context{
				Namespace: "test-namespace",
			},
		},
		{
			name:           "Cluster context name is set, return context",
			clusterContext: "test-context",
			config: &clientcmdapi.Config{
				Contexts: map[string]*clientcmdapi.Context{
					"test-context": {
						Namespace: "test-namespace",
					},
				},
			},
			expectedCtx: &clientcmdapi.Context{
				Namespace: "test-namespace",
			},
		},
		{
			name:           "Cluster context name is set, but context is not available, return nil context",
			clusterContext: "test-context",
			config: &clientcmdapi.Config{
				Contexts: map[string]*clientcmdapi.Context{
					"other-context": {
						Namespace: "other-namespace",
					},
				},
			},
			expectedCtx: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterContextName = tt.clusterContext
			SetClientConfigAPI(tt.config)
			actualCtx := GetContext()
			assert.Equal(t, tt.expectedCtx, actualCtx)
		})
	}
}

func tearDown() {
	SetClientConfigAPI(nil)
	SetClusterContextName("")
	SetConfigClusterServerName("")
	SetK8SGitServerVersion("")
}
