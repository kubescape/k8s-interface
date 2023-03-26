package v1

import (
	"context"
	"fmt"
	"os"
	"strings"

	container "cloud.google.com/go/container/apiv1"
	"github.com/kubescape/k8s-interface/k8sinterface"
	"golang.org/x/oauth2/google"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type IGKESupport interface {
	GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error)
	GetName(clusterDescribe *containerpb.Cluster) string
	GetProject(cluster string) (string, error)
	GetRegion(cluster string) (string, error)
	GetContextName(cluster string) string
}
type GKESupport struct {
}

var (
	KS_GKE_PROJECT_ENV_VAR = "KS_GKE_PROJECT"
)

func NewGKESupport() *GKESupport {
	return &GKESupport{}
}

func (gkeSupport *GKESupport) GetRegion(cluster string) (string, error) {
	region, present := os.LookupEnv(KS_CLOUD_REGION_ENV_VAR)
	if present {
		return region, nil
	}
	parsedName := strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return "", fmt.Errorf("failed to parse region name from cluster name: '%s'", cluster)
	}
	region = parsedName[2]
	return region, nil
}

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

func (gkeSupport *GKESupport) GetName(clusterDescribe *containerpb.Cluster) string {
	return clusterDescribe.Name
}

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

func (gkeSupport *GKESupport) GetContextName(cluster string) string {

	parsedName := strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return ""
	}
	clusterName := parsedName[3]
	if clusterName != "" {
		return clusterName
	}
	cluster = k8sinterface.GetContextName()
	parsedName = strings.Split(cluster, "_")
	if len(parsedName) < 3 {
		return ""
	}
	return parsedName[3]
}
