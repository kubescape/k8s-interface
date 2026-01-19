package cloudmetadata

import (
	"context"
	"fmt"
	"strings"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	corev1 "k8s.io/api/core/v1"
)

const (
	// Providers moved to armoapi-go/armotypes
	TestMode = "testmode"
)

// GetCloudMetadata retrieves cloud metadata for a given node
func GetCloudMetadata(ctx context.Context, node *corev1.Node, nodeName string) (*armotypes.CloudMetadata, error) {
	metadata := &armotypes.CloudMetadata{
		Hostname: node.Name,
		HostType: armotypes.HostTypeKubernetes,
	}

	// Determine provider and extract metadata
	providerID := node.Spec.ProviderID
	switch {
	case strings.HasPrefix(providerID, "aws://"):
		metadata.Provider = armotypes.ProviderAws
		metadata = extractAWSMetadata(node, metadata)
	case strings.HasPrefix(providerID, "gce://"):
		metadata.Provider = armotypes.ProviderGcp
		metadata = extractGCPMetadata(node, metadata)
	case strings.HasPrefix(providerID, "azure://"):
		metadata.Provider = armotypes.ProviderAzure
		metadata = extractAzureMetadata(node, metadata)
	case strings.HasPrefix(providerID, "digitalocean://"):
		metadata.Provider = armotypes.ProviderDigitalOcean
		metadata = extractDigitalOceanMetadata(node, metadata)
	case strings.HasPrefix(providerID, "openstack://"):
		metadata.Provider = armotypes.ProviderOpenStack
		metadata = extractOpenstackMetadata(node, metadata)
	case strings.HasPrefix(providerID, "vsphere://"):
		metadata.Provider = armotypes.ProviderVMware
		metadata = extractVMwareMetadata(node, metadata)
	case strings.HasPrefix(providerID, "alicloud://"):
		metadata.Provider = armotypes.ProviderAlibaba
		metadata = extractAlibabaMetadata(node, metadata)
	case strings.HasPrefix(providerID, "ibm://"):
		metadata.Provider = armotypes.ProviderIBM
		metadata = extractIBMMetadata(node, metadata)
	case strings.HasPrefix(providerID, "oci://"):
		metadata.Provider = armotypes.ProviderOracle
		metadata = extractOracleMetadata(node, metadata)
	case strings.HasPrefix(providerID, "linode://"):
		metadata.Provider = armotypes.ProviderLinode
		metadata = extractLinodeMetadata(node, metadata)
	case strings.HasPrefix(providerID, "scaleway://"):
		metadata.Provider = armotypes.ProviderScaleway
		metadata = extractScalewayMetadata(node, metadata)
	case strings.HasPrefix(providerID, "vultr://"):
		metadata.Provider = armotypes.ProviderVultr
		metadata = extractVultrMetadata(node, metadata)
	case strings.HasPrefix(providerID, "hcloud://"):
		metadata.Provider = armotypes.ProviderHetzner
		metadata = extractHetznerMetadata(node, metadata)
	case strings.HasPrefix(providerID, "equinixmetal://"):
		metadata.Provider = armotypes.ProviderEquinixMetal
		metadata = extractEquinixMetalMetadata(node, metadata)
	case strings.HasPrefix(providerID, "exoscale://"):
		metadata.Provider = armotypes.ProviderExoscale
		metadata = extractExoscaleMetadata(node, metadata)
	default:
		metadata.Provider = armotypes.ProviderOther
		if v := ctx.Value(TestMode); v != nil {
			logger.L().Ctx(ctx).Warning("Test mode: unknown cloud provider for node %s: %s", helpers.String("nodeName", nodeName), helpers.String("providerID", providerID))
		} else {
			return nil, fmt.Errorf("unknown cloud provider for node %s: %s", nodeName, providerID)
		}
	}

	// Extract common metadata from node status
	for _, addr := range node.Status.Addresses {
		switch addr.Type {
		case "InternalIP":
			metadata.PrivateIP = addr.Address
		case "ExternalIP":
			metadata.PublicIP = addr.Address
		case "Hostname":
			metadata.Hostname = addr.Address
		}
	}

	return metadata, nil
}

func extractAWSMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	metadata.HostType = armotypes.HostTypeEc2
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: aws:///us-west-2a/i-1234567890abcdef0
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract account ID from annotations if available
	if accountID, ok := node.Annotations["eks.amazonaws.com/account-id"]; ok {
		metadata.AccountID = accountID
		if node.Labels["eks.amazonaws.com/compute-type"] == "fargate" {
			metadata.HostType = armotypes.HostTypeEksFargate
		} else {
			metadata.HostType = armotypes.HostTypeEksEc2
		}
	} else {
		// Extract account ID from metadata service if available
		client := ec2metadata.New(session.Must(session.NewSession()))
		if client.Available() {
			identity, err := client.GetInstanceIdentityDocument()
			if err == nil {
				metadata.AccountID = identity.AccountID
			}
		}
	}

	return metadata
}

func extractGCPMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	metadata.HostType = armotypes.HostTypeGce
	// Extract from labels
	metadata.InstanceType = node.Labels["beta.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	if _, ok := node.Labels["cloud.google.com/gke-nodepool"]; ok {
		if node.Labels["cloud.google.com/gke-provisioning"] == "autopilot" || node.Labels["cloud.google.com/gke-autopilot"] == "true" {
			metadata.HostType = armotypes.HostTypeAutopilot
		} else {
			metadata.HostType = armotypes.HostTypeGke
		}
	}

	// Extract project and instance ID from provider ID
	// Format: gce://project-name/zone/instance-name
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 3 {
		metadata.AccountID = parts[2] // project name
		metadata.InstanceID = parts[len(parts)-1]
	}

	return metadata
}

func extractAzureMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	metadata.HostType = armotypes.HostTypeAzureVm
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	if _, ok := node.Labels["kubernetes.azure.com/cluster"]; ok || node.Labels["kubernetes.azure.com/agentpool"] != "" {
		metadata.HostType = armotypes.HostTypeAks
	}

	// Extract subscription ID and resource info from provider ID
	// Format: azure:///subscriptions/<id>/resourceGroups/<name>/providers/Microsoft.Compute/virtualMachineScaleSets/<name>
	if parts := strings.Split(node.Spec.ProviderID, "/"); len(parts) > 3 {
		for i, part := range parts {
			if part == "subscriptions" && i+1 < len(parts) {
				metadata.AccountID = parts[i+1]
			}
			if part == "virtualMachineScaleSets" && i+1 < len(parts) {
				metadata.InstanceID = parts[i+1]
			}
		}
	}

	return metadata
}

func extractDigitalOceanMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	metadata.HostType = armotypes.HostTypeDroplet
	// Extract from labels
	metadata.InstanceType = node.Labels["beta.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	if node.Labels["doks.digitalocean.com/managed"] == "true" || node.Labels["doks.digitalocean.com/node-id"] != "" {
		metadata.HostType = armotypes.HostTypeDoks
	}

	// Extract droplet ID from provider ID
	// Format: digitalocean:///droplet-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	return metadata
}

func extractOpenstackMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: openstack:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract project ID if available
	if projectID, ok := node.Labels["project.openstack.org/project-id"]; ok {
		metadata.AccountID = projectID
	}

	return metadata
}

func extractVMwareMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract VM UUID from provider ID
	// Format: vsphere:///vm-uuid
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract datacenter info if available
	if dc, ok := node.Labels["vsphere.kubernetes.io/datacenter"]; ok {
		metadata.Region = dc
	}

	return metadata
}

func extractAlibabaMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: alicloud:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract account ID if available
	if accountID, ok := node.Labels["alibabacloud.com/account-id"]; ok {
		metadata.AccountID = accountID
	}

	return metadata
}

func extractIBMMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: ibm:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract account ID if available
	if accountID, ok := node.Labels["ibm-cloud.kubernetes.io/account-id"]; ok {
		metadata.AccountID = accountID
	}

	return metadata
}

func extractOracleMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract OCID from provider ID
	// Format: oci:///ocid
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract compartment ID if available
	if compartmentID, ok := node.Labels["oci.oraclecloud.com/compartment-id"]; ok {
		metadata.AccountID = compartmentID
	}

	return metadata
}

func extractLinodeMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Use Linode-specific region label if available
	if liNodeRegion, ok := node.Labels["topology.linode.com/region"]; ok && metadata.Region == "" {
		metadata.Region = liNodeRegion
	}

	// Extract Linode ID from provider ID
	// Format: linode:///linode-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Check for Linode-specific private IP annotation
	if privateIP, ok := node.Annotations["node.k8s.linode.com/private-ip"]; ok && metadata.PrivateIP == "" {
		metadata.PrivateIP = privateIP
	}

	// Check for Linode-specific hostname label
	if hostname, ok := node.Labels["kubernetes.io/hostname"]; ok && metadata.Hostname == "" {
		metadata.Hostname = hostname
	}

	return metadata
}

func extractScalewayMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: scaleway:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract organization ID if available
	if orgID, ok := node.Labels["scaleway.com/organization-id"]; ok {
		metadata.AccountID = orgID
	}

	return metadata
}

func extractVultrMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: vultr:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	return metadata
}

func extractHetznerMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract server ID from provider ID
	// Format: hcloud:///server-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract project ID if available
	if projectID, ok := node.Labels["hcloud.hetzner.cloud/project-id"]; ok {
		metadata.AccountID = projectID
	}

	return metadata
}

func extractEquinixMetalMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract device ID from provider ID
	// Format: equinixmetal:///device-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract project ID if available
	if projectID, ok := node.Labels["metal.equinix.com/project-id"]; ok {
		metadata.AccountID = projectID
	}

	return metadata
}

func extractExoscaleMetadata(node *corev1.Node, metadata *armotypes.CloudMetadata) *armotypes.CloudMetadata {
	// Extract from labels
	metadata.InstanceType = node.Labels["node.kubernetes.io/instance-type"]
	metadata.Region = node.Labels["topology.kubernetes.io/region"]
	metadata.Zone = node.Labels["topology.kubernetes.io/zone"]

	// Extract instance ID from provider ID
	// Format: exoscale:///instance-id
	parts := strings.Split(node.Spec.ProviderID, "/")
	if len(parts) > 0 {
		metadata.InstanceID = parts[len(parts)-1]
	}

	// Extract organization ID if available
	if orgID, ok := node.Labels["exoscale.com/organization-id"]; ok {
		metadata.AccountID = orgID
	}

	return metadata
}
