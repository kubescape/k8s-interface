package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

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
	GetPolicyVersion(region string) (*ListPolicyVersion, error)
}

type EKSSupport struct {
}

const (
	awsauthconfigmap = "aws-auth"

	// Bounds AWS SDK and IAM/ECR calls so an air-gapped or otherwise
	// unreachable AWS control plane cannot stall the scan loop.
	eksCallTimeout            = 5 * time.Second
	eksRBACEnumerationTimeout = 30 * time.Second
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

type PolicyVersionDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

type ListPolicyVersion struct {
	PolicyVersion map[string]*PolicyVersionDocument `json:"policiesDocuments"`
}

func NewEKSSupport() *EKSSupport {
	return &EKSSupport{}
}

func (eksSupport *EKSSupport) GetClusterDescribe(cluster string, region string) (*eks.DescribeClusterOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), eksCallTimeout)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := eks.NewFromConfig(awsConfig)
	input := &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	}

	result, err := svc.DescribeCluster(ctx, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (eksSupport *EKSSupport) GetName(describe *eks.DescribeClusterOutput) string {
	return *describe.Cluster.Name
}

func (eksSupport *EKSSupport) GetRegion(cluster string) (string, error) {
	region, present := os.LookupEnv(KS_CLOUD_REGION_ENV_VAR)
	if present && region != "" {
		return region, nil
	}

	region, present = os.LookupEnv("AWS_REGION")
	if present && region != "" {
		return region, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), eksCallTimeout)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err == nil && awsConfig.Region != "" {
		return awsConfig.Region, nil
	}

	if parsed, err := arn.Parse(cluster); err == nil && parsed.Region != "" {
		return parsed.Region, nil
	}

	splittedClusterContext := strings.Split(cluster, "-")
	if len(splittedClusterContext) >= 6 {
		return strings.Join(splittedClusterContext[3:6], "-"), nil
	}

	return "", fmt.Errorf("failed to get region: tried environment variables (KS_CLOUD_REGION, AWS_REGION), AWS config, and cluster name parsing")
}

func (eksSupport *EKSSupport) GetContextName(cluster string) string {
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

	splittedCluster = strings.Split(cluster, ":")
	if len(splittedCluster) > 5 {
		clusterName := splittedCluster[len(splittedCluster)-1]
		clusterNameFiltered := strings.Replace(clusterName, "cluster-", "", 1)
		if clusterName != clusterNameFiltered {
			return clusterNameFiltered
		}
	}

	splittedCluster = strings.Split(k8sinterface.GetContextName(), "/")
	if len(splittedCluster) > 1 {
		return splittedCluster[len(splittedCluster)-1]
	}

	if cluster != "" {
		return cluster
	}

	return ""
}

func (EKSSupport *EKSSupport) GetEKSCfgMap(kapi *k8sinterface.KubernetesApi, namespace string) (*v1.ConfigMap, error) {
	var authData awsAuth

	ctx, cancel := context.WithTimeout(context.Background(), eksCallTimeout)
	defer cancel()

	eksCfgMap, err := kapi.KubernetesClient.CoreV1().ConfigMaps(namespace).Get(ctx, awsauthconfigmap, metav1.GetOptions{})
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

func (eksSupport *EKSSupport) GetDescribeRepositories(region string) (*ecr.DescribeRepositoriesOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), eksRBACEnumerationTimeout)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := ecr.NewFromConfig(awsConfig)
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(100),
	}

	result, err := svc.DescribeRepositories(ctx, input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (eksSupport *EKSSupport) GetListEntitiesForPolicies(region string) (*ListEntitiesForPolicies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), eksRBACEnumerationTimeout)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	svc := iam.NewFromConfig(awsConfig)
	input := &iam.ListPoliciesInput{}

	result, err := listPoliciesWithPagination(ctx, svc, input)
	if err != nil {
		return nil, err
	}
	allEntitiesForPolicies := map[string]*iam.ListEntitiesForPolicyOutput{}
	for _, policy := range result {
		inp := &iam.ListEntitiesForPolicyInput{
			PolicyArn: policy.Arn,
		}
		entitiesForPolicy, err := svc.ListEntitiesForPolicy(ctx, inp)
		if err != nil {
			return nil, err
		}
		allEntitiesForPolicies[*policy.Arn] = entitiesForPolicy
	}
	return &ListEntitiesForPolicies{EntitiesForPolicies: allEntitiesForPolicies}, nil
}

func (eksSupport *EKSSupport) GetPolicyVersion(region string) (*ListPolicyVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), eksRBACEnumerationTimeout)
	defer cancel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error: fail to load AWS SDK default %v", err)
	}
	awsConfig.Region = region
	svc := iam.NewFromConfig(awsConfig)

	input := &iam.ListPoliciesInput{}
	result, err := listPoliciesWithPagination(ctx, svc, input)
	if err != nil {
		return nil, fmt.Errorf("error: fail to list policies: %v", err)
	}

	policyVersionContents := map[string]*PolicyVersionDocument{}
	for _, policy := range result {
		policyVersionInput := &iam.GetPolicyVersionInput{
			PolicyArn: policy.Arn,
			VersionId: policy.DefaultVersionId,
		}
		policyVersionContent, err := svc.GetPolicyVersion(ctx, policyVersionInput)
		if err != nil {
			return nil, fmt.Errorf("error: fail to get policy version: %v", err)
		}
		policyVersionDocument, err := url.QueryUnescape(*policyVersionContent.PolicyVersion.Document)
		if err != nil {
			return nil, fmt.Errorf("error: fail to decode Document field: %v", err)
		}
		pDocument := PolicyVersionDocument{}
		json.Unmarshal([]byte(policyVersionDocument), &pDocument)

		policyVersionContents[*policy.Arn] = &pDocument
	}
	return &ListPolicyVersion{PolicyVersion: policyVersionContents}, nil
}

func listPoliciesWithPagination(ctx context.Context, svc *iam.Client, input *iam.ListPoliciesInput) ([]types.Policy, error) {
	paginator := iam.NewListPoliciesPaginator(svc, input)

	var policiesList []types.Policy
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error: fail to list policies: %v", err)
		}
		for _, policy := range output.Policies {
			policiesList = append(policiesList, policy)
		}
	}
	return policiesList, nil
}
