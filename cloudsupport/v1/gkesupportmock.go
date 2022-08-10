package v1

import (
	"encoding/json"
	"strings"

	"github.com/kubescape/k8s-interface/cloudsupport/mockobjects"
	"github.com/kubescape/k8s-interface/k8sinterface"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func NewGKESupportMock() *GKESupportMock {
	return &GKESupportMock{}
}

type GKESupportMock struct {
}

// Get descriptive info about cluster running in GKE.
func (gkeSupportM *GKESupportMock) GetClusterDescribe(cluster string, region string, project string) (*containerpb.Cluster, error) {
	c := &containerpb.Cluster{}
	err := json.Unmarshal([]byte(mockobjects.GkeDescriptor), c)
	return c, err
}

func (gkeSupportM *GKESupportMock) GetName(clusterDescribe *containerpb.Cluster) string {
	return clusterDescribe.Name
}

func (gkeSupportM *GKESupportMock) GetProject(cluster string) (string, error) {
	return "", nil
}

func (gkeSupportM *GKESupportMock) GetRegion(cluster string) (string, error) {
	return "", nil
}

func (gkeSupportM *GKESupportMock) GetContextName(cluster string) string {
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
