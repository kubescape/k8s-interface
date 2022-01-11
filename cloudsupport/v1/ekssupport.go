package v1

import (
	"fmt"
	"strings"

	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

type IEKSSupport interface {
	GetClusterDescribe(currContext string) (*eks.DescribeClusterOutput, error)
	GetName(*eks.DescribeClusterOutput) string
}

type EKSSupport struct {
}

func NewEKSSupport() *EKSSupport {
	return &EKSSupport{}
}

// Get descriptive info about cluster running in EKS.
func (eksSupport *EKSSupport) GetClusterDescribe(currContext string) (*eks.DescribeClusterOutput, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	splittedClusterContext := strings.Split(currContext, ".")
	if len(splittedClusterContext) < 2 {
		return nil, fmt.Errorf("error: failed to get region")
	}
	region := splittedClusterContext[1]

	// Configure cluster name and region for request
	svc := eks.New(s, &aws.Config{Region: aws.String(region)})
	input := &eks.DescribeClusterInput{
		Name: aws.String(k8sinterface.GetClusterName()),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getName get cluster name from describe
func (eksSupport *EKSSupport) GetName(describe *eks.DescribeClusterOutput) string {
	return *describe.Cluster.Name
}
