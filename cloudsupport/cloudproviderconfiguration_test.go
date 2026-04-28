package cloudsupport

import (
	"errors"
	"os"
	"testing"

	cloudsupportv1 "github.com/kubescape/k8s-interface/cloudsupport/v1"
)

func TestKSOfflineShortCircuitsCloudDescribe(t *testing.T) {
	t.Setenv(KS_OFFLINE_ENV_VAR, "true")

	entryPoints := map[string]func(string, string) (interface{}, error){
		"GetDescriptiveInfoFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetDescriptiveInfoFromCloudProvider(c, p)
		},
		"GetDescribeRepositoriesFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetDescribeRepositoriesFromCloudProvider(c, p)
		},
		"GetListEntitiesForPoliciesFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetListEntitiesForPoliciesFromCloudProvider(c, p)
		},
		"GetPolicyVersionFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetPolicyVersionFromCloudProvider(c, p)
		},
	}

	for name, fn := range entryPoints {
		t.Run(name, func(t *testing.T) {
			got, err := fn("any-cluster", cloudsupportv1.AKS)
			if got != nil {
				t.Fatalf("%s: expected nil result when offline, got %v", name, got)
			}
			if !errors.Is(err, ErrCloudDescribeUnavailable) {
				t.Fatalf("%s: expected ErrCloudDescribeUnavailable, got %v", name, err)
			}
		})
	}
}

// AKS path returns ErrCloudDescribeUnavailable (rather than a bare error) when
// AZURE_SUBSCRIPTION_ID / AZURE_RESOURCE_GROUP are not set. This is the
// failure mode real air-gapped users hit before any network call is attempted,
// and it's what lets the kubescape scan loop classify the failure as
// non-fatal.
func TestAKSMissingCredsWrapsSentinel(t *testing.T) {
	os.Unsetenv(KS_OFFLINE_ENV_VAR)
	os.Unsetenv(cloudsupportv1.AZURE_SUBSCRIPTION_ID_ENV_VAR)
	os.Unsetenv(cloudsupportv1.AZURE_RESOURCE_GROUP_ENV_VAR)

	cases := map[string]func(string, string) (interface{}, error){
		"GetDescriptiveInfoFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetDescriptiveInfoFromCloudProvider(c, p)
		},
		"GetListEntitiesForPoliciesFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetListEntitiesForPoliciesFromCloudProvider(c, p)
		},
		"GetPolicyVersionFromCloudProvider": func(c, p string) (interface{}, error) {
			return GetPolicyVersionFromCloudProvider(c, p)
		},
	}

	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := fn("any-cluster", cloudsupportv1.AKS)
			if !errors.Is(err, ErrCloudDescribeUnavailable) {
				t.Fatalf("%s: expected ErrCloudDescribeUnavailable, got %v", name, err)
			}
		})
	}
}
