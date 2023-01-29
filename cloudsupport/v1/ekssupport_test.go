package v1

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContextName(t *testing.T) {
	mockname1 := "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster"
	eksSupport := NewEKSSupport()
	name := eksSupport.GetContextName(mockname1)
	assert.Equal(t, "test-cluster", name)
	region, err := eksSupport.GetRegion(mockname1)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	mockname2 := "arn:aws:eks:eu-north-1:123456789:cluster/test-cluster"
	splittedCluster := strings.Split(mockname2, "/")
	name = splittedCluster[len(splittedCluster)-1]
	assert.Equal(t, "test-cluster", name)
	region, err = eksSupport.GetRegion(mockname2)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

}

func TestGetRegion(t *testing.T) {
	mockname1 := "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster"
	eksSupport := NewEKSSupport()
	region, err := eksSupport.GetRegion(mockname1)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	mockname2 := "arn-aws-eks-eu-west-2-XXXXXXXXXXXX-cluster-Yiscah-test-g2am5"
	region, err = eksSupport.GetRegion(mockname2)
	assert.NoError(t, err)
	assert.Equal(t, "eu-west-2", region)
}
