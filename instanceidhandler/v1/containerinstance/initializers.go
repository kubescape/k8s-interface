package containerinstance

import (
	"fmt"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
)

// GenerateInstanceIDFromString generates instance ID from string.
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
//
// For Job workloads (Kind == "Job"), this function uses a best-effort heuristic to recover AlternateName
// and TemplateHash when the Job name ends in a suffix produced by the CronJob pod template hasher (rand.SafeEncodeString).
// Because the serialized string format does not preserve provenance metadata, standalone literal Jobs whose names
// happen to end in matching suffixes (such as "backup-db" or "backup-4") will also have those suffixes parsed as
// template hashes.
func GenerateInstanceIDFromString(input string) (*InstanceID, error) {

	instanceID := &InstanceID{}

	if fields := RegexFormatted.FindStringSubmatch(input); fields != nil {
		instanceID.ApiVersion = fields[1]
		instanceID.Namespace = fields[2]
		instanceID.Kind = fields[3]
		instanceID.Name = fields[4]
		instanceID.InstanceType = fields[5]
		instanceID.ContainerName = fields[6]
	} else if fields := RegexNoContainer.FindStringSubmatch(input); fields != nil {
		instanceID.ApiVersion = fields[1]
		instanceID.Namespace = fields[2]
		instanceID.Kind = fields[3]
		instanceID.Name = fields[4]
	} else {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	// Contextual recovery contract:
	// For Job instance IDs (which may represent Pods owned by CronJob child Jobs where GetStringFormatted()
	// serialized AlternateName into the name field), infer AlternateName and TemplateHash if the name suffix
	// matches the producer's template hash domain.
	// Ordinary CronJob workloads (Kind == "CronJob") and standard Jobs without a valid template hash suffix
	// remain literal with empty AlternateName and TemplateHash.
	// As documented above, standalone literal Jobs whose names end in a matching producer suffix (e.g. "backup-db"
	// or "backup-4") are inferred as template hashes under this best-effort tradeoff.
	if instanceID.Kind == "Job" {
		s := strings.Split(instanceID.Name, "-")
		if len(s) > 1 && helpers.IsTemplateHash(s[len(s)-1]) {
			instanceID.AlternateName = instanceID.Name
			instanceID.TemplateHash = s[len(s)-1]
		}
	}

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	// Check if the input string is valid
	if instanceID.GetStringFormatted() != input && instanceID.GetStringNoContainer() != input {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	return instanceID, nil
}
