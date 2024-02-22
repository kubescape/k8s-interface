package v1

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContextName(t *testing.T) {
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
}
