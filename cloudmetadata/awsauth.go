package cloudmetadata

import (
	"fmt"
	"strings"

	apitypes "github.com/armosec/armoapi-go/armotypes"
	corev1 "k8s.io/api/core/v1"
)

// EnrichCloudMetadataFromAWSAuthConfigMap enriches cloud metadata account ID from aws-auth ConfigMap
func EnrichCloudMetadataFromAWSAuthConfigMap(metadata *apitypes.CloudMetadata, cm *corev1.ConfigMap) error {
	if metadata == nil || metadata.Provider != ProviderAWS || metadata.AccountID != "" {
		return nil
	}

	if cm.Name != "aws-auth" || cm.Namespace != "kube-system" {
		return fmt.Errorf("expected aws-auth ConfigMap in kube-system namespace, got %s in %s", cm.Name, cm.Namespace)
	}

	// Get the mapRoles data from ConfigMap
	mapRolesData, ok := cm.Data["mapRoles"]
	if !ok {
		return fmt.Errorf("mapRoles data not found in aws-auth ConfigMap")
	}

	lines := strings.Split(mapRolesData, "\n")

	for _, line := range lines {
		if strings.Contains(line, "rolearn") {
			line = strings.TrimLeft(line, " -")
			line = strings.TrimPrefix(line, "rolearn:")
			line = strings.TrimSpace(line)

			// Parse the ARN to extract account ID
			// ARN format: arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME
			parts := strings.Split(line, ":")
			if len(parts) >= 5 {
				metadata.AccountID = parts[4]
				return nil
			}
		}
	}

	return fmt.Errorf("no valid role ARN found in aws-auth ConfigMap")
}
