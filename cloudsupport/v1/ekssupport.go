package v1

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"

	//"github.com/aws/aws-sdk-go-v2/aws/session"
	"github.com/armosec/k8s-interface/k8sinterface"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type IEKSSupport interface {
	GetClusterDescribe(currContext string, region string) (*eks.DescribeClusterOutput, error)
	GetName(*eks.DescribeClusterOutput) string
	GetRegion(cluster string) (string, error)
	GetCluster(cluster string) string
}

type EKSSupport struct {
}

func NewEKSSupport() *EKSSupport {
	return &EKSSupport{}
}

// Get descriptive info about cluster running in EKS.

func (eksSupport *EKSSupport) GetClusterDescribe(cluster string, region string) (*eks.DescribeClusterOutput, error) {
	// Configure cluster name and region for request
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := eks.NewFromConfig(awsConfig)
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
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

func (eksSupport *EKSSupport) GetRegion(cluster string) (string, error) {
	region, present := os.LookupEnv(KS_CLOUD_REGION_ENV_VAR)
	if present {
		return region, nil
	}
	splittedClusterContext := strings.Split(cluster, ".")

	if len(splittedClusterContext) < 2 {
		splittedClusterContext := strings.Split(cluster, ":")
		if len(splittedClusterContext) < 4 {
			return "", fmt.Errorf("failed to get region")
		}
		region = splittedClusterContext[3]
	} else {
		region = splittedClusterContext[1]
	}
	return region, nil
}

func (eksSupport *EKSSupport) GetCluster(cluster string) string {
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
