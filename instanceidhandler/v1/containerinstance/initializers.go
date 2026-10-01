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

	// Handle CronJob child Jobs where the name contains a template hash
	if instanceID.Kind == "Job" || instanceID.Kind == "CronJob" {
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
