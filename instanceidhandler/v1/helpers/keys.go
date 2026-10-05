package helpers

import (
	"math"
	"regexp"
	"strconv"

	"github.com/kubescape/k8s-interface/k8sinterface"
)

type InstanceType string

// metadata keys
const (
	metadataPrefix                     = "kubescape.io"
	ApiGroupMetadataKey                = metadataPrefix + "/workload-api-group"
	ApiVersionMetadataKey              = metadataPrefix + "/workload-api-version"
	CloudAccountIdentifierMetadataKey  = metadataPrefix + "/cloud-account-identifier"
	ClusterMetadataKey                 = metadataPrefix + "/cluster"
	ClusterRoleBindingNameMetadataKey  = metadataPrefix + "/clusterrolebinding-name"
	ClusterRoleNameMetadataKey         = metadataPrefix + "/clusterrole-name"
	CompletionMetadataKey              = metadataPrefix + "/completion"
	ContainerNameMetadataKey           = metadataPrefix + "/workload-container-name"
	ContainerTypeMetadataKey           = metadataPrefix + "/workload-container-type"
	ContextMetadataKey                 = metadataPrefix + "/context"
	EphemeralContainerNameMetadataKey  = metadataPrefix + "/workload-ephemeral-container-name" // DEPRECATED - use ContainerNameMetadataKey and ContainerTypeMetadataKey
	HostIDMetadataKey                  = metadataPrefix + "/host-id"
	HostTypeMetadataKey                = metadataPrefix + "/host-type"
	ImageIDMetadataKey                 = metadataPrefix + "/image-id"
	ImageNameMetadataKey               = metadataPrefix + "/image-name"
	ImageTagMetadataKey                = metadataPrefix + "/image-tag"
	InitContainerNameMetadataKey       = metadataPrefix + "/workload-init-container-name" // DEPRECATED - use ContainerNameMetadataKey and ContainerTypeMetadataKey
	InstanceIDMetadataKey              = metadataPrefix + "/instance-id"
	LearningPeriodMetadataKey          = metadataPrefix + "/learning-period"
	OtelSpanIDMetadataKey              = metadataPrefix + "/otel-span-id"
	OtelTraceparentMetadataKey         = metadataPrefix + "/otel-traceparent"
	ManagedByMetadataKey               = metadataPrefix + "/managed-by"
	PreviousReportTimestampMetadataKey = metadataPrefix + "/previous-report-timestamp"
	RbacResourceMetadataKey            = metadataPrefix + "/rbac-resource"
	RegionMetadataKey                  = metadataPrefix + "/region"
	RelatedKindMetadataKey             = metadataPrefix + "/workload-kind"
	RelatedNameMetadataKey             = metadataPrefix + "/workload-name"
	RelatedNamespaceMetadataKey        = metadataPrefix + "/workload-namespace"
	ReportSeriesIdMetadataKey          = metadataPrefix + "/report-series-id"
	ReportTimestampMetadataKey         = metadataPrefix + "/report-timestamp"
	ResourceSizeMetadataKey            = metadataPrefix + "/resource-size"
	ResourceVersionMetadataKey         = metadataPrefix + "/workload-resource-version"
	RoleBindingNameMetadataKey         = metadataPrefix + "/rolebinding-name"
	RoleBindingNamespaceMetadataKey    = metadataPrefix + "/rolebinding-namespace"
	RoleNameMetadataKey                = metadataPrefix + "/role-name"
	RoleNamespaceMetadataKey           = metadataPrefix + "/role-namespace"
	ScanIdMetadataKey                  = metadataPrefix + "/scan-id"
	StatusMetadataKey                  = metadataPrefix + "/status"
	SyncChecksumMetadataKey            = metadataPrefix + "/sync-checksum"
	TemplateHashKey                    = metadataPrefix + "/instance-template-hash"
	TierMetadataKey                    = metadataPrefix + "/tier"
	ToolVersionMetadataKey             = metadataPrefix + "/tool-version"
	UserDefinedProfileMetadataKey      = metadataPrefix + "/user-defined-profile" // should be used as a label!
	WlidMetadataKey                    = metadataPrefix + "/wlid"
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

// sbom metadata keys and values
const (
	ArtifactTypeMetadataKey = metadataPrefix + "/sbom-type"
	ContainerArtifactType   = "container"
	HostArtifactType        = "host"
	ImageArtifactType       = "image"
	NodeArtifactType        = "node"
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
	Initializing      = "initializing"
	Learning          = "ready"
	Completed         = "completed"
	Incomplete        = "incomplete"
	Unauthorize       = "unauthorize"
	MissingRuntime    = "missing-runtime"
	TooLarge          = "too-large"
	Failed            = "failed" // container exited with a non-zero code
	UnsupportedSchema = "unsupported-schema"
)

// Completion
const (
	Partial = "partial"
	Full    = "complete"
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

// TemplateHashRegex matches the output domain of rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32())).
// In k8s.io/apimachinery/pkg/util/rand, SafeEncodeString encodes decimal digits '0'-'9' using alphanums = "bcdfghjklmnpqrstvwxz2456789".
// Because (int('0')+d) % 27 maps decimal digits 0..9 strictly to '4','5','6','7','8','9','b','c','d','f',
// the producer domain consists exclusively of characters [4-9bcdf], with length 1 to 10 (decimal representation of uint32).
var TemplateHashRegex = regexp.MustCompile(`^[4-9bcdf]{1,10}$`)

// IsUnixTimeInMinutes checks if a string represents a Unix timestamp in minutes.
// In Kubernetes, CronJob controller appends the schedule timestamp (minutes since epoch) to the child Job name.
func IsUnixTimeInMinutes(s string) bool {
	if i, err := strconv.Atoi(s); err == nil {
		return i > 0 && i < math.MaxInt64/60
	}
	return false
}

// IsTemplateHash returns true if s matches the template hash domain produced by
// rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32())).
//
// Best-Effort Contextual Recovery & Tradeoff Contract:
// In Kubernetes, CronJob child Jobs are originally named "<cronjob>-<timestamp>" where timestamp
// is the scheduled Unix time in minutes (e.g. "backup-28677846"). When Kubescape generates an InstanceID
// for a Pod under such a Job, it computes a template hash via SafeEncodeString and sets AlternateName
// to "<cronjob>-<templateHash>" (e.g. "backup-8698448884" or "kubevuln-scheduler-b449cf78f").
// When serialized via GetStringFormatted(), AlternateName replaces the timestamp-bearing name in the name field.
//
// Because the legacy serialized string format does not retain provenance metadata:
//  1. Literal names or timestamps containing non-encoder characters (e.g. '0', '1', '2', '3' present in
//     standard Unix timestamps like "28677846", or 'x' in literal names like "backup-xxxxxxxx")
//     are rejected and preserved literally.
//  2. Alphanumeric suffixes containing 'b', 'c', 'd', 'f' cannot be Unix timestamps (which are purely decimal).
//     However, standalone literal Jobs whose names happen to end in matching suffixes (e.g. "backup-db")
//     will be inferred as having template hashes by this best-effort heuristic.
//  3. Purely numeric suffixes (digits '4'-'9', such as "8698448884" or "backup-4") that represent valid
//     canonical uint32 encodings are recovered as template hashes, accepting the tradeoff that literal
//     Jobs with identical suffixes are indistinguishable without external provenance.
func IsTemplateHash(s string) bool {
	if !TemplateHashRegex.MatchString(s) {
		return false
	}
	// Validate canonical uint32 encoding:
	// '4' decodes to '0'. Canonical fmt.Sprint(uint32) has no leading zero unless s == "4" (0).
	if len(s) > 1 && s[0] == '4' {
		return false
	}
	// Decode back to decimal digits to ensure the value does not exceed math.MaxUint32 (4294967295)
	if len(s) == 10 {
		decoded := make([]byte, 10)
		for i := 0; i < 10; i++ {
			c := s[i]
			switch {
			case c >= '4' && c <= '9':
				decoded[i] = c - 4
			case c == 'b':
				decoded[i] = '6'
			case c == 'c':
				decoded[i] = '7'
			case c == 'd':
				decoded[i] = '8'
			case c == 'f':
				decoded[i] = '9'
			}
		}
		if _, err := strconv.ParseUint(string(decoded), 10, 32); err != nil {
			return false
		}
	}
	return true
}
