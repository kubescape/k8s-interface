package containerinstance

import (
	"fmt"
	"strings"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
)

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
func GenerateInstanceIDFromString(input string) (*InstanceID, error) {

	instanceID := &InstanceID{}

	// TODO add case for CronJobs here, or deprecate

	// Split the input string by the field separator "/"
	fields := strings.Split(input, helpers.StringFormatSeparator)
	if len(fields) != 5 && len(fields) != 6 {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	i := 0
	instanceID.ApiVersion = strings.TrimPrefix(fields[0], helpers.PrefixApiVersion)

	// if the apiVersion has a group, e.g. apps/v1
	if len(fields) == 6 {
		instanceID.ApiVersion += helpers.StringFormatSeparator + fields[1]
		i += 1
	}

	instanceID.Namespace = strings.TrimPrefix(fields[1+i], helpers.PrefixNamespace)
	instanceID.Kind = strings.TrimPrefix(fields[2+i], helpers.PrefixKind)
	instanceID.Name = strings.TrimPrefix(fields[3+i], helpers.PrefixName)
	instanceID.InstanceType, instanceID.ContainerName, _ = strings.Cut(fields[4+i], "Name-")

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	// Check if the input string is valid
	if instanceID.GetStringFormatted() != input {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	return instanceID, nil
}
