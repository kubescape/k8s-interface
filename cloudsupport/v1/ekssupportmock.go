package v1

import (
	"encoding/json"
	"strings"

	"github.com/armosec/k8s-interface/cloudsupport/mockobjects"
	"github.com/armosec/k8s-interface/k8sinterface"
	"github.com/aws/aws-sdk-go-v2/service/eks"
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
			return splittedCluster[0]
		}
	}
	splittedCluster := strings.Split(k8sinterface.GetContextName(), ".")
	if len(splittedCluster) > 1 {
		return splittedCluster[0]
	}
	return ""
}
