package cloudsupport

import (
	"context"
	"fmt"
	"strings"

	container "cloud.google.com/go/container/apiv1"
	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type GKESupport struct {
}

func newGKESupport() *GKESupport {
	return &GKESupport{}
}

// Get descriptive info about cluster running in GKE.
func (gkeSupport *GKESupport) getClusterDescribe() (*containerpb.Cluster, error) {
	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	parsedName := strings.Split(k8sinterface.GetClusterName(), "_")
	if len(parsedName) < 3 {
		return nil, fmt.Errorf("error: failed to parse cluster name")
	}
	clusterName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", parsedName[1], parsedName[2], parsedName[3])
	req := &containerpb.GetClusterRequest{
		Name: clusterName,
	}
	result, err := c.GetCluster(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (gkeSupport *GKESupport) getName(clusterDescribe *containerpb.Cluster) string {
	return clusterDescribe.Name
}
