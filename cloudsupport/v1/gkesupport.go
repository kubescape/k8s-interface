package v1

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	container "cloud.google.com/go/container/apiv1"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type IGKESupport interface {
	GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error)
	GetName(clusterDescribe *containerpb.Cluster) string
	GetProject(cluster string) (string, error)
	GetRegion(cluster string) (string, error)
	GetContextName(cluster string) string
	GetDescribeRepositories(project string, region string) ([]*artifactregistrypb.Repository, error)
	GetListEntitiesForPolicies(project string) (*cloudresourcemanager.Policy, error)
	GetIAMMappings(project string) (map[string]string, map[string]string, error)
}
type GKESupport struct {
}

var (
	KS_GKE_PROJECT_ENV_VAR = "KS_GKE_PROJECT"
)

const (
	// gkeRBACEnumerationTimeout is the budget for paginated enumeration calls which take longer
	gkeRBACEnumerationTimeout = 30 * time.Second
)

func NewGKESupport() *GKESupport {
	return &GKESupport{}
}

// GetRegion extracts the GCP region from the cluster string or environment variable
func (gkeSupport *GKESupport) GetRegion(cluster string) (string, error) {
	region, present := os.LookupEnv(KS_CLOUD_REGION_ENV_VAR)
	if present && strings.TrimSpace(region) != "" {
		return region, nil
	}
	parsedName := strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return "", fmt.Errorf("failed to parse region from cluster name: '%s'", cluster)
	}
	region = parsedName[2]
	return region, nil
}

// GetDescribeRepositories returns a list of GCP Artifact Registries in the given project and region
func (gkeSupport *GKESupport) GetDescribeRepositories(project string, region string) ([]*artifactregistrypb.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gkeRBACEnumerationTimeout)
	defer cancel()

	client, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Parent format: projects/PROJECT_ID/locations/LOCATION_ID
	// Note: We deliberately limit the scope to a single region (the cluster's region) in this PR.
	// We also normalize zonal clusters to their region (e.g. us-central1-c -> us-central1)
	// since a GCP zone is always a region plus a single-letter suffix.
	normalizedRegion := region
	parts := strings.Split(region, "-")
	if len(parts) == 3 && len(parts[2]) == 1 {
		normalizedRegion = strings.Join(parts[:2], "-")
	}
	req := &artifactregistrypb.ListRepositoriesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, normalizedRegion),
	}

	it := client.ListRepositories(ctx, req)
	var repositories []*artifactregistrypb.Repository

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, resp)
	}

	return repositories, nil
}

// GetListEntitiesForPolicies returns a list of IAM role bindings in the given project
func (gkeSupport *GKESupport) GetListEntitiesForPolicies(project string) (*cloudresourcemanager.Policy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gkeRBACEnumerationTimeout)
	defer cancel()

	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	req := &cloudresourcemanager.GetIamPolicyRequest{
		Options: &cloudresourcemanager.GetPolicyOptions{
			RequestedPolicyVersion: 3,
		},
	}
	policy, err := crmService.Projects.GetIamPolicy(project, req).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// GetProject extracts the GCP project from the cluster string or environment variable
func (gkeSupport *GKESupport) GetProject(cluster string) (string, error) {
	project, present := os.LookupEnv(KS_GKE_PROJECT_ENV_VAR)
	if present {
		return project, nil
	}
	parsedName := strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return "", fmt.Errorf("failed to parse project name from cluster name: '%s'", cluster)
	}
	project = parsedName[1]
	return project, nil
}

// Get descriptive info about cluster running in GKE.
func (gkeSupport *GKESupport) GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error) {
	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	clusterName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, region, cluster)
	req := &containerpb.GetClusterRequest{
		Name: clusterName,
	}
	result, err := c.GetCluster(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetName returns the name of the GKE cluster
func (gkeSupport *GKESupport) GetName(clusterDescribe *containerpb.Cluster) string {
	return clusterDescribe.Name
}

// GetAuthorizationKey retrieves the GCP access token
func (gkeSupport *GKESupport) GetAuthorizationKey() (string, error) {
	ctx := context.Background()

	token, err := google.DefaultTokenSource(ctx, nil...)
	if err != nil {
		return "", fmt.Errorf("failed to find creds: %w", err)
	}
	t, err := token.Token()
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	return t.AccessToken, nil
}

// GetContextName extracts the context name from the cluster string
func (gkeSupport *GKESupport) GetContextName(cluster string) string {

	parsedName := strings.Split(cluster, "_")
	if len(parsedName) != 4 {
		return ""
	}
	clusterName := parsedName[3]
	if clusterName != "" {
		return clusterName
	}
	cluster = k8sinterface.GetContextName()
	parsedName = strings.Split(cluster, "_")
	if len(parsedName) != 4 {
		return ""
	}
	return parsedName[3]
}

// GetIAMMappings returns iam-roles and service accounts
func (gkeSupport *GKESupport) GetIAMMappings(project string) (map[string]string, map[string]string, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, iam.CloudPlatformScope)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Google Cloud client: %w", err)
	}

	iamService, err := iam.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create IAM service client: %w", err)
	}

	roleMappings := make(map[string]string)
	saMappings := make(map[string]string)

	roleIterator, err := iamService.Projects.Roles.List("projects/" + project).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve roles: %w", err)
	}
	for {
		for _, role := range roleIterator.Roles {
			roleMappings[role.Name] = role.Title
		}

		if roleIterator.NextPageToken == "" {
			break
		}

		roleIterator, err = iamService.Projects.Roles.List("projects/" + project).PageToken(roleIterator.NextPageToken).Do()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to retrieve roles: %w", err)
		}
	}

	saIterator, err := iamService.Projects.ServiceAccounts.List("projects/" + project).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve service accounts: %w", err)
	}
	for {

		for _, sa := range saIterator.Accounts {
			saMappings[sa.Name] = sa.Name
		}

		if saIterator.NextPageToken == "" {
			break
		}

		saIterator, err = iamService.Projects.ServiceAccounts.List("projects/" + project).PageToken(saIterator.NextPageToken).Do()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to retrieve service accounts: %w", err)
		}
	}

	return roleMappings, saMappings, nil
}
