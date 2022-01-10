package v1

import (
	"context"
	"fmt"
	"strings"

	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/aws/aws-sdk-go-v2/aws"

	//"github.com/aws/aws-sdk-go-v2/aws/session"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"k8s.io/client-go/tools/clientcmd/api"
)

type IEKSSupport interface {
	GetClusterDescribe(currContext *api.Context) (*eks.DescribeClusterOutput, error)
	GetName(*eks.DescribeClusterOutput) string
}

type EKSSupport struct {
}

func NewEKSSupport() *EKSSupport {
	return &EKSSupport{}
}

// Get descriptive info about cluster running in EKS.
func (eksSupport *EKSSupport) GetClusterDescribe(currContext *api.Context) (*eks.DescribeClusterOutput, error) {
	splittedClusterContext := strings.Split(k8sinterface.GetCurrentContext().Cluster, ".")
	if len(splittedClusterContext) < 2 {
		return nil, fmt.Errorf("error: failed to get region")
	}
	region := splittedClusterContext[1]

	// Configure cluster name and region for request
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := eks.NewFromConfig(awsConfig)
	input := &eks.DescribeClusterInput{
		Name: aws.String(k8sinterface.GetClusterName()),
	}

	result, err := svc.DescribeCluster(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getName get cluster name from describe
func (eksSupport *EKSSupport) GetName(describe *eks.DescribeClusterOutput) string {
	return *describe.Cluster.Name
}
