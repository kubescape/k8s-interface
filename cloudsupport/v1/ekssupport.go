package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"github.com/aws/aws-sdk-go-v2/aws/session"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kubescape/k8s-interface/k8sinterface"
)

type IEKSSupport interface {
	GetClusterDescribe(currContext string, region string) (*eks.DescribeClusterOutput, error)
	GetName(*eks.DescribeClusterOutput) string
	GetRegion(cluster string) (string, error)
	GetContextName(cluster string) string
	GetDescribeRepositories(region string) (*ecr.DescribeRepositoriesOutput, error)
	GetListEntitiesForPolicies(region string) (*ListEntitiesForPolicies, error)
}

type EKSSupport struct {
}

const (
	awsauthconfigmap = "aws-auth"
)

type awsAuth struct {
	MapRoles []*mappedRoles `json:"mapRoles"`
	MapUsers []*mappedUsers `json:"mapUsers"`
}

type mappedRoles struct {
	RoleArn  string   `json:"rolearn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups,omitempty"`
}

type mappedUsers struct {
	UserArn  string   `json:"userarn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups,omitempty"`
}

type ListEntitiesForPolicies struct {
	EntitiesForPolicies map[string]*iam.ListEntitiesForPolicyOutput `json:"rolesPolicies"`
}

// NewEKSSupport returns EKSSupport type
func NewEKSSupport() *EKSSupport {
	return &EKSSupport{}
}

// GetClusterDescribe returns the descriptive info about the cluster running in EKS.
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

// GetName returns the name of the eks cluster
func (eksSupport *EKSSupport) GetName(describe *eks.DescribeClusterOutput) string {

	//getName get cluster name from describe
	return *describe.Cluster.Name
}

// GetRegion returns the region in which eks cluster is running.
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

// Context can be in one of 2 ways:
// 1. arn:aws:eks:<region>:<id>:cluster/<cluster_name> --> Usually this will be in context
// 2. arn:aws:eks:<region>:<id>:cluster-<cluster_name> --> Usually we will get 'cluster' param like this
func (eksSupport *EKSSupport) GetContextName(cluster string) string {
	if cluster != "" {
		splittedCluster := strings.Split(cluster, ".")
		if len(splittedCluster) > 1 {
			return splittedCluster[0]
		}
	}
	// Try from context
	splittedCluster := strings.Split(k8sinterface.GetContextName(), ".")
	if len(splittedCluster) > 1 {
		return splittedCluster[0]
	}

	splittedCluster = strings.Split(cluster, ":")
	if len(splittedCluster) > 5 {
		// arn:aws:eks:<region>:<id>:cluster-<cluster_name> -> <cluster_name>
		clusterName := splittedCluster[len(splittedCluster)-1]
		clusterNameFiltered := strings.Replace(clusterName, "cluster-", "", 1)
		if clusterName != clusterNameFiltered {
			return clusterNameFiltered
		}
	}

	// Try from context
	splittedCluster = strings.Split(k8sinterface.GetContextName(), "/")
	if len(splittedCluster) > 1 {
		// arn:aws:eks:<region>:<id>:cluster/<cluster_name> -> <cluster_name>
		return splittedCluster[len(splittedCluster)-1]
	}
	return ""
}

// GetEKSCfgMap returns the ConfigMap containing mappings of iam-roles/groups or iam-users/groups
func (EKSSupport *EKSSupport) GetEKSCfgMap(kapi *k8sinterface.KubernetesApi, namespace string) (*v1.ConfigMap, error) {

	var authData awsAuth

	eksCfgMap, err := kapi.KubernetesClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), awsauthconfigmap, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	if mapRoles, ok := eksCfgMap.Data["mapRoles"]; ok {

		if err := json.Unmarshal([]byte(mapRoles), &authData.MapRoles); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("'mapRoles' is missing from the EKS config object")
	}

	if mapUsers, ok := eksCfgMap.Data["mapUsers"]; ok {

		if err := json.Unmarshal([]byte(mapUsers), &authData.MapUsers); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("'mapUsers' is missing from the EKS config object")
	}

	return eksCfgMap, nil

}

// GetDescribeRepositories returns the descriptive info about the repositories in EKS.
func (eksSupport *EKSSupport) GetDescribeRepositories(region string) (*ecr.DescribeRepositoriesOutput, error) {
	// Configure region for request
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := ecr.NewFromConfig(awsConfig)
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(100),
	}

	result, err := svc.DescribeRepositories(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetListEntitiesForPolicies returns the list of roles in EKS.
func (eksSupport *EKSSupport) GetListEntitiesForPolicies(region string) (*ListEntitiesForPolicies, error) {
	// Configure region for request
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	svc := iam.NewFromConfig(awsConfig)
	input := &iam.ListPoliciesInput{}

	//TODO - Add pagination
	result, err := svc.ListPolicies(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	allEntitiesForPolicies := map[string]*iam.ListEntitiesForPolicyOutput{}
	for _, policy := range result.Policies {
		inp := &iam.ListEntitiesForPolicyInput{
			PolicyArn: policy.Arn,
		}
		entitiesForPolicy, err := svc.ListEntitiesForPolicy(context.TODO(), inp)
		if err != nil {
			return nil, err
		}
		allEntitiesForPolicies[*policy.Arn] = entitiesForPolicy
	}
	return &ListEntitiesForPolicies{EntitiesForPolicies: allEntitiesForPolicies}, nil
}
