package v1

import (
	"testing"

	"github.com/armosec/k8s-interface/cloudsupport/apis"
	"github.com/armosec/k8s-interface/k8sinterface"
	"github.com/stretchr/testify/assert"
)

func TestGetClusterDescribeGKE(t *testing.T) {
	g := NewGKESupportMock()
	des, err := GetClusterDescribeGKE(g)
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderDescribeKind, des.GetKind())
	assert.Equal(t, "cloud.google.com/v1/Describe/kubescape-demo-01", des.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionGKE, Version), des.GetApiVersion())
	assert.Equal(t, "kubescape-demo-01", des.GetName())
	assert.Equal(t, 34, len(des.GetData()))
}

func TestGetClusterDescribeEKS(t *testing.T) {
	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, nil)
	assert.NoError(t, err)
	assert.Equal(t, apis.CloudProviderDescribeKind, des.GetKind())
	assert.Equal(t, "eks.amazonaws.com/v1/Describe/ca-terraform-eks-dev-stage", des.GetID())
	assert.Equal(t, k8sinterface.JoinGroupVersion(apis.ApiVersionEKS, Version), des.GetApiVersion())
	assert.Equal(t, "ca-terraform-eks-dev-stage", des.GetName())
	assert.Equal(t, 1, len(des.GetData()))
}

func TestNewDescriptiveInfoFromCloudProvider(t *testing.T) {
	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, nil)
	assert.NoError(t, err)

	assert.True(t, IsTypeDescriptiveInfoFromCloudProvider(des.GetObject()))
	d := NewDescriptiveInfoFromCloudProvider(des.GetObject())
	assert.NotNil(t, d)

	assert.Equal(t, d.GetID(), des.GetID())

}
func TestSetObject(t *testing.T) {
	g := NewEKSSupportMock()
	des, err := GetClusterDescribeEKS(g, nil)
	assert.NoError(t, err)

	d := CloudProviderDescribe{}
	d.SetObject(des.GetObject())

	assert.Equal(t, d.GetID(), des.GetID())
	assertMap(t, d.GetObject(), des.GetObject())
}

func assertMap(t *testing.T, expected, actual map[string]interface{}) {
	for k0, v0 := range expected {
		if v1, ok := actual[k0]; ok {
			switch f := v1.(type) {
			case string:
				assert.Equal(t, v0.(string), f)
			case map[string]interface{}:
				assertMap(t, v0.(map[string]interface{}), f)
			}
		}
	}
}
