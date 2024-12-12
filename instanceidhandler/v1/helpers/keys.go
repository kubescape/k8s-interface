package helpers

import "github.com/kubescape/k8s-interface/k8sinterface"

type InstanceType string

// metadata keys
const (
	metadataPrefix                    = "kubescape.io"
	ApiGroupMetadataKey               = metadataPrefix + "/workload-api-group"
	ApiVersionMetadataKey             = metadataPrefix + "/workload-api-version"
	ClusterRoleBindingNameMetadataKey = metadataPrefix + "/clusterrolebinding-name"
	ClusterRoleNameMetadataKey        = metadataPrefix + "/clusterrole-name"
	CompletionMetadataKey             = metadataPrefix + "/completion"
	ContainerNameMetadataKey          = metadataPrefix + "/workload-container-name"
	ContextMetadataKey                = metadataPrefix + "/context"
	EphemeralContainerNameMetadataKey = metadataPrefix + "/workload-ephemeral-container-name"
	ImageIDMetadataKey                = metadataPrefix + "/image-id"
	ImageNameMetadataKey              = metadataPrefix + "/image-name"
	ImageTagMetadataKey               = metadataPrefix + "/image-tag"
	InitContainerNameMetadataKey      = metadataPrefix + "/workload-init-container-name"
	InstanceIDMetadataKey             = metadataPrefix + "/instance-id"
	KindMetadataKey                   = metadataPrefix + "/workload-kind"
	ManagedByMetadataKey              = metadataPrefix + "/managed-by"
	NameMetadataKey                   = metadataPrefix + "/workload-name"
	NamespaceMetadataKey              = metadataPrefix + "/workload-namespace"
	RbacResourceMetadataKey           = metadataPrefix + "/rbac-resource"
	ResourceSizeMetadataKey           = metadataPrefix + "/resource-size"
	ResourceVersionMetadataKey        = metadataPrefix + "/workload-resource-version"
	RoleBindingNameMetadataKey        = metadataPrefix + "/rolebinding-name"
	RoleBindingNamespaceMetadataKey   = metadataPrefix + "/rolebinding-namespace"
	RoleNameMetadataKey               = metadataPrefix + "/role-name"
	RoleNamespaceMetadataKey          = metadataPrefix + "/role-namespace"
	ScanIdMetadataKey                 = metadataPrefix + "/scan-id"
	StatusMetadataKey                 = metadataPrefix + "/status"
	SyncChecksumMetadataKey           = metadataPrefix + "/sync-checksum"
	TemplateHashKey                   = metadataPrefix + "/instance-template-hash"
	TierMetadataKey                   = metadataPrefix + "/tier"
	WlidMetadataKey                   = metadataPrefix + "/wlid"
)

// metadata values
const (
	ContextMetadataKeyFiltered    = "filtered"
	ContextMetadataKeyNonFiltered = "non-filtered"
)

// application profile metadata
const (
	ManagedByUserValue            = "User"
	UserApplicationProfilePrefix  = "ug-"
	UserNetworkNeighborhoodPrefix = "ug-"
)

// string format: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/...
const (
	StringFormatSeparator = "/"
	PrefixApiVersion      = "apiVersion-"
	PrefixNamespace       = "namespace-"
	PrefixKind            = "kind-"
	PrefixName            = "name-"
)

// Statuses
const (
	Initializing   = "initializing"
	Ready          = "ready"
	Completed      = "completed"
	Incomplete     = "incomplete"
	Unauthorize    = "unauthorize"
	MissingRuntime = "missing-runtime"
	TooLarge       = "too-large"
)

// Completion
const (
	Partial  = "partial"
	Complete = "complete"
)

// Tier values
const (
	CoreTier = "core"
)

func IgnoreOwnerReference(ownerKind string) bool {
	if ownerKind == "Node" {
		return true
	}
	if _, e := k8sinterface.GetGroupVersionResource(ownerKind); e != nil {
		return true
	}
	return false
}
