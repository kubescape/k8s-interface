package v1

import (
	"context"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	"golang.org/x/oauth2/google"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type IGKESupport interface {
	GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error)
	GetName(clusterDescribe *containerpb.Cluster) string
}
type GKESupport struct {
}

func NewGKESupport() *GKESupport {
	return &GKESupport{}
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
		fmt.Printf("GetAuthorizationKey: DefaultTokenSource failed with error: %v\n", err)
		return "", fmt.Errorf("failed to find creds")
	}
	t, err := token.Token()
	if err != nil {
		fmt.Printf("GetAuthorizationKey: DefaultTokenSource failed with error: %v\n", err)
		return "", fmt.Errorf("failed to find creds")
	}
	fmt.Printf("GetAuthorizationKey: t.AccessToken %v\n", t.AccessToken)
	return t.AccessToken, nil
}
