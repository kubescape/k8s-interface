package helpers

import "github.com/kubescape/k8s-interface/k8sinterface"

type InstanceType string

// metadata keys
const (
	metadataPrefix                    = "kubescape.io"
	ApiGroupMetadataKey               = metadataPrefix + "/workload-api-group"
	ApiVersionMetadataKey             = metadataPrefix + "/workload-api-version"
	ContainerNameMetadataKey          = metadataPrefix + "/workload-container-name"
	InitContainerNameMetadataKey      = metadataPrefix + "/workload-init-container-name"
	EphemeralContainerNameMetadataKey = metadataPrefix + "/workload-ephemeral-container-name"
	ImageNameMetadataKey              = metadataPrefix + "/image-name"
	ImageTagMetadataKey               = metadataPrefix + "/image-tag"
	ImageIDMetadataKey                = metadataPrefix + "/image-id"
	InstanceIDMetadataKey             = metadataPrefix + "/instance-id"
	TemplateHashKey                   = metadataPrefix + "/instance-template-hash"
	KindMetadataKey                   = metadataPrefix + "/workload-kind"
	NameMetadataKey                   = metadataPrefix + "/workload-name"
	NamespaceMetadataKey              = metadataPrefix + "/workload-namespace"
	ResourceVersionMetadataKey        = metadataPrefix + "/workload-resource-version"
	StatusMetadataKey                 = metadataPrefix + "/status"
	WlidMetadataKey                   = metadataPrefix + "/wlid"
	ContextMetadataKey                = metadataPrefix + "/context"
	RbacResourceMetadataKey           = metadataPrefix + "/rbac-resource"
	RoleNameMetadataKey               = metadataPrefix + "/role-name"
	RoleNamespaceMetadataKey          = metadataPrefix + "/role-namespace"
	RoleBindingNameMetadataKey        = metadataPrefix + "/rolebinding-name"
	RoleBindingNamespaceMetadataKey   = metadataPrefix + "/rolebinding-namespace"
	ClusterRoleNameMetadataKey        = metadataPrefix + "/clusterrole-name"
	ClusterRoleBindingNameMetadataKey = metadataPrefix + "/clusterrolebinding-name"
	ResourceSizeMetadataKey           = metadataPrefix + "/resource-size"
	CompletionMetadataKey             = metadataPrefix + "/completion"
	TierMetadataKey                   = metadataPrefix + "/tier"
	ManagedByMetadataKey              = metadataPrefix + "/managed-by"
	ApplicationProfileCommand         = metadataPrefix + "/application-profile-command"
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
	AddCommand                    = "add"
	DeleteCommand                 = "delete"
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
