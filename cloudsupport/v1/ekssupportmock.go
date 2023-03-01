package v1

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kubescape/k8s-interface/cloudsupport/mockobjects"
	"github.com/kubescape/k8s-interface/k8sinterface"
)

func NewEKSSupportMock() *EKSSupportMock {
	return &EKSSupportMock{}
}

type EKSSupportMock struct {
}

// Get descriptive info about cluster running in EKS.
func (eksSupportM *EKSSupportMock) GetClusterDescribe(currContext string, region string) (*eks.DescribeClusterOutput, error) {
	describeClusterOutput := &eks.DescribeClusterOutput{}
	err := json.Unmarshal([]byte(mockobjects.EksDescriptor), describeClusterOutput)
	return describeClusterOutput, err
}

// GetDescribeRepositories
func (eksSupportM *EKSSupportMock) GetDescribeRepositories(region string) (*ecr.DescribeRepositoriesOutput, error) {
	describeRepositoriesOutput := &ecr.DescribeRepositoriesOutput{}
	err := json.Unmarshal([]byte(mockobjects.EksDescribeRepositories), describeRepositoriesOutput)
	return describeRepositoriesOutput, err
}

// GetListEntitiesForPolicies
func (eksSupportM *EKSSupportMock) GetListEntitiesForPolicies(region string) (*ListEntitiesForPolicies, error) {
	listEntitiesForPoliciesOutput := &ListEntitiesForPolicies{}
	err := json.Unmarshal([]byte(mockobjects.EksListEntitiesForPolicies), listEntitiesForPoliciesOutput)
	return listEntitiesForPoliciesOutput, err
}

// GetPolicyVersion
func (eksSupportM *EKSSupportMock) GetPolicyVersion(region string) (*ListPolicyVersion, error) {
	policyVersionContent := &ListPolicyVersion{}
	err := json.Unmarshal([]byte(mockobjects.EksGetPolicyVersion), policyVersionContent)
	return policyVersionContent, err
}

// getName get cluster name from describe
func (eksSupportM *EKSSupportMock) GetName(describe *eks.DescribeClusterOutput) string {
	return *describe.Cluster.Name
}

func (eksSupportM *EKSSupportMock) GetRegion(cluster string) (string, error) {
	return "", nil
}

func (eksSupport *EKSSupportMock) GetContextName(cluster string) string {
	if cluster != "" {
		splittedCluster := strings.Split(cluster, ".")
		if len(splittedCluster) > 1 {
			return splittedCluster[len(splittedCluster)-1]
		}
	}
	splittedCluster := strings.Split(k8sinterface.GetContextName(), ".")
	if len(splittedCluster) > 1 {
		return splittedCluster[len(splittedCluster)-1]
	}
	splittedCluster = strings.Split(cluster, ":")
	if len(splittedCluster) > 1 {
		return splittedCluster[len(splittedCluster)-1]
	}
	splittedCluster = strings.Split(k8sinterface.GetContextName(), ":")
	if len(splittedCluster) > 1 {
		return splittedCluster[len(splittedCluster)-1]
	}
	return ""
}
