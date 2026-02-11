package v1

import (
	"os"
	"strings"
	"testing"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/stretchr/testify/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetContextName(t *testing.T) {
	defer tearDown()

	// Test ARN context names
	mockname1 := "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster"
	eksSupport := NewEKSSupport()
	name := eksSupport.GetContextName(mockname1)
	assert.Equal(t, "test-cluster", name)
	region, err := eksSupport.GetRegion(mockname1)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	mockname2 := "arn:aws:eks:eu-north-1:123456789:cluster/test-cluster"
	splittedCluster := strings.Split(mockname2, "/")
	name = splittedCluster[len(splittedCluster)-1]
	assert.Equal(t, "test-cluster", name)
	region, err = eksSupport.GetRegion(mockname2)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	// Test non-ARN context names (i.e., cluster name is context name)
	tests := []struct {
		name                string
		config              *clientcmdapi.Config
		connected           bool
		cluster             string
		expectedContextName string
	}{
		{
			name: "Config is empty, return empty string",
			config: &clientcmdapi.Config{
				CurrentContext: "d34db33f",
				Clusters: map[string]*clientcmdapi.Cluster{
					"d34db33f": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           true,
			cluster:             "my-cluster",
			expectedContextName: "",
		},
		{
			name: "Context name is cluster name, connected",
			config: &clientcmdapi.Config{
				CurrentContext: "my-cluster",
				Clusters: map[string]*clientcmdapi.Cluster{
					"my-cluster": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           true,
			cluster:             "my-cluster",
			expectedContextName: "my-cluster",
		},
		{
			name: "Context name is cluster name, not connected",
			config: &clientcmdapi.Config{
				CurrentContext: "my-cluster",
				Clusters: map[string]*clientcmdapi.Cluster{
					"my-cluster": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           false,
			cluster:             "my-cluster",
			expectedContextName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sinterface.SetConnectedToCluster(tt.connected)
			k8sinterface.SetClientConfigAPI(tt.config)
			actualContextName := eksSupport.GetContextName(tt.cluster)
			assert.Equal(t, tt.expectedContextName, actualContextName)
		})
	}

}

func TestGetRegion(t *testing.T) {
	tests := []struct {
		name           string
		cluster        string
		envRegion      string
		expectedRegion string
		expectErr      bool
	}{
		{
			name:           "Region is extracted from cluster name 1",
			cluster:        "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster",
			expectedRegion: "eu-north-1",
			expectErr:      false,
		},
		{
			name:           "Region is extracted from cluster name 2",
			cluster:        "arn-aws-eks-eu-west-2-XXXXXXXXXXXX-cluster-Yiscah-test-g2am5",
			expectedRegion: "eu-west-2",
			expectErr:      false,
		},
		{
			name:           "Region is present in environment variable",
			envRegion:      "us-west-2",
			expectedRegion: "us-west-2",
			expectErr:      false,
		},
		{
			name:           "Region is extracted from cluster name with ':' separator",
			cluster:        "cluster:us-west-2:eks",
			expectedRegion: "us-west-2",
			expectErr:      true,
		},
		{
			name:      "Failed to get region",
			cluster:   "cluster",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envRegion != "" {
				os.Setenv(KS_CLOUD_REGION_ENV_VAR, tt.envRegion)
				defer os.Unsetenv(KS_CLOUD_REGION_ENV_VAR)
			}

			eksSupport := &EKSSupport{}
			region, err := eksSupport.GetRegion(tt.cluster)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRegion, region)
			}
		})
	}

	t.Run("AWS_REGION environment variable is checked", func(t *testing.T) {
		awsRegion := "ap-southeast-2"
		os.Setenv("AWS_REGION", awsRegion)
		defer os.Unsetenv("AWS_REGION")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, awsRegion, region)
	})

	t.Run("KS_CLOUD_REGION takes precedence over AWS_REGION", func(t *testing.T) {
		ksRegion := "us-west-2"
		awsRegion := "eu-west-1"
		os.Setenv(KS_CLOUD_REGION_ENV_VAR, ksRegion)
		os.Setenv("AWS_REGION", awsRegion)
		defer os.Unsetenv(KS_CLOUD_REGION_ENV_VAR)
		defer os.Unsetenv("AWS_REGION")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, ksRegion, region)
	})
}
