package containerinstance

import (
	"fmt"
)

// GenerateInstanceIDFromString generates instance ID from string
// The string format is: apiVersion-<apiVersion>/namespace-<namespace>/kind-<kind>/name-<name>/containerName-<containerName>
func GenerateInstanceIDFromString(input string) (*InstanceID, error) {

	instanceID := &InstanceID{}

	// TODO add case for CronJobs here, or deprecate

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

	if err := validateInstanceID(instanceID); err != nil {
		return nil, err
	}

	// Check if the input string is valid
	if instanceID.GetStringFormatted() != input && instanceID.GetStringNoContainer() != input {
		return nil, fmt.Errorf("invalid format: %s", input)
	}

	return instanceID, nil
}
