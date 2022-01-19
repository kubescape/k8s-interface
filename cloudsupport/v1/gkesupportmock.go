package v1

import (
	"encoding/json"

	"github.com/armosec/k8s-interface/cloudsupport/mockobjects"
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
