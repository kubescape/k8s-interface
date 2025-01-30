package cloudmetadata

import (
	"testing"

	apitypes "github.com/armosec/armoapi-go/armotypes"
	"github.com/kubescape/k8s-interface/utils"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

	_ "embed"
)

//go:embed testdata/awsauth.json
var awsauth []byte

func TestEnrichCloudMetadataFromAWSAuthConfigMap(t *testing.T) {
	cm, err := utils.ConvertUnstructuredToRuntimeObject(utils.MustConvertRawToUnstructured(awsauth))
	if err != nil {
		t.Fatalf("ConvertUnstructuredToRuntimeObject() error = %v", err)
	}
	configmapObj := cm.(*v1.ConfigMap)
	metadata := &apitypes.CloudMetadata{
		Provider: ProviderAWS,
	}
	err = EnrichCloudMetadataFromAWSAuthConfigMap(metadata, configmapObj)
	assert.NoError(t, err)
	assert.Equal(t, "012345678912", metadata.AccountID)
}
