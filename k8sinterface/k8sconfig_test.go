package k8sinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetK8sConfigClusterServerName(t *testing.T) {
	tests := []struct {
		name           string
		config         *clientcmdapi.Config
		expectedServer string
	}{
		{
			name:           "Config is nil, return default server name",
			config:         nil,
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
			actualServer := GetK8sConfigClusterServerName(tt.config)
			assert.Equal(t, tt.expectedServer, actualServer)
		})
	}
}
