package v1

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/stretchr/testify/assert"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()

	previousValue, hadPreviousValue := os.LookupEnv(key)
	if value == "" {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
	} else {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("set %s: %v", key, err)
		}
	}

	t.Cleanup(func() {
		var err error
		if hadPreviousValue {
			err = os.Setenv(key, previousValue)
		} else {
			err = os.Unsetenv(key)
		}

		if err != nil {
			t.Fatalf("restore %s: %v", key, err)
		}
	})
}

func isolateAWSRegionSources(t *testing.T) {
	t.Helper()

	configPath := filepath.Join(t.TempDir(), "config")
	credentialsPath := filepath.Join(t.TempDir(), "credentials")

	if err := os.WriteFile(configPath, nil, 0o600); err != nil {
		t.Fatalf("write %s: %v", configPath, err)
	}
	if err := os.WriteFile(credentialsPath, nil, 0o600); err != nil {
		t.Fatalf("write %s: %v", credentialsPath, err)
	}

	setEnv(t, KS_CLOUD_REGION_ENV_VAR, "")
	setEnv(t, "AWS_REGION", "")
	setEnv(t, "AWS_DEFAULT_REGION", "")
	setEnv(t, "AWS_PROFILE", "")
	setEnv(t, "AWS_EC2_METADATA_DISABLED", "true")
	setEnv(t, "AWS_CONFIG_FILE", configPath)
	setEnv(t, "AWS_SHARED_CREDENTIALS_FILE", credentialsPath)
}

func TestGetContextName(t *testing.T) {
	isolateAWSRegionSources(t)

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
			isolateAWSRegionSources(t)

			if tt.envRegion != "" {
				setEnv(t, KS_CLOUD_REGION_ENV_VAR, tt.envRegion)
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
		isolateAWSRegionSources(t)

		awsRegion := "ap-southeast-2"
		setEnv(t, "AWS_REGION", awsRegion)

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, awsRegion, region)
	})

	t.Run("KS_CLOUD_REGION takes precedence over AWS_REGION", func(t *testing.T) {
		isolateAWSRegionSources(t)

		ksRegion := "us-west-2"
		awsRegion := "eu-west-1"
		setEnv(t, KS_CLOUD_REGION_ENV_VAR, ksRegion)
		setEnv(t, "AWS_REGION", awsRegion)

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, ksRegion, region)
	})

	t.Run("KS_CLOUD_REGION takes precedence over cluster name parsing", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, KS_CLOUD_REGION_ENV_VAR, "us-west-2")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("arn:aws:eks:eu-north-1:123456789:cluster/test-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "us-west-2", region)
	})

	t.Run("Cluster name parsing takes precedence over AWS default config", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_DEFAULT_REGION", "us-east-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("arn:aws:eks:eu-north-1:123456789:cluster/test-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "eu-north-1", region)
	})

	t.Run("Cluster name parsing takes precedence over AWS_REGION", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_REGION", "us-east-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("arn:aws:eks:eu-north-1:123456789:cluster/test-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "eu-north-1", region)
	})

	t.Run("Partitioned ARN parsing takes precedence over AWS_REGION", func(t *testing.T) {
		tests := []struct {
			name    string
			cluster string
			region  string
		}{
			{
				name:    "GovCloud standard ARN",
				cluster: "arn:aws-us-gov:eks:us-gov-west-1:123456789012:cluster/test-cluster",
				region:  "us-gov-west-1",
			},
			{
				name:    "China standard ARN",
				cluster: "arn:aws-cn:eks:cn-north-1:123456789012:cluster/test-cluster",
				region:  "cn-north-1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				isolateAWSRegionSources(t)
				setEnv(t, "AWS_REGION", "us-east-1")

				eksSupport := &EKSSupport{}
				region, err := eksSupport.GetRegion(tt.cluster)

				assert.NoError(t, err)
				assert.Equal(t, tt.region, region)
			})
		}
	})

	t.Run("Dashed ARN parsing supports GovCloud regions", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_REGION", "us-east-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("arn-aws-us-gov-eks-us-gov-west-1-123456789012-cluster-test-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "us-gov-west-1", region)
	})

	t.Run("Valid dotted cluster name takes precedence over AWS fallbacks", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_REGION", "us-east-1")
		setEnv(t, "AWS_DEFAULT_REGION", "ap-southeast-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("cluster.us-west-2.eksctl.io")

		assert.NoError(t, err)
		assert.Equal(t, "us-west-2", region)
	})

	t.Run("Invalid dotted cluster name falls back to AWS_REGION", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_REGION", "us-west-2")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("foo.bar")

		assert.NoError(t, err)
		assert.Equal(t, "us-west-2", region)
	})

	t.Run("Invalid parsed region falls back to AWS default config", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_DEFAULT_REGION", "ap-southeast-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("arn:aws:eks:not-a-region:123456789:cluster/test-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "ap-southeast-1", region)
	})

	t.Run("AWS default config is used after cluster parsing fails", func(t *testing.T) {
		isolateAWSRegionSources(t)
		setEnv(t, "AWS_DEFAULT_REGION", "ap-southeast-1")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, "ap-southeast-1", region)
	})
}
