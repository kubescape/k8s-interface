package ecsinstance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler/v1/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestECSInstanceID_TaskDefinition(t *testing.T) {
	tests := []struct {
		name           string
		clusterName    string
		taskDefinition string
		revision       string
		containerName  string
	}{
		{
			name:           "basic task definition",
			clusterName:    "my-cluster",
			taskDefinition: "my-app",
			revision:       "1",
			containerName:  "app-container",
		},
		{
			name:           "task definition with long name",
			clusterName:    "production-cluster",
			taskDefinition: "very-long-task-definition-name-for-testing",
			revision:       "42",
			containerName:  "main",
		},
		{
			name:           "task definition with hyphens",
			clusterName:    "test-cluster",
			taskDefinition: "my-web-service-prod",
			revision:       "23",
			containerName:  "nginx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instanceID := NewTaskDefinitionInstanceID(tt.clusterName, tt.taskDefinition, tt.revision, tt.containerName)

			assert.Equal(t, ECSDefaultApiVersion, instanceID.ApiVersion)
			assert.Equal(t, tt.clusterName, instanceID.ClusterName)
			assert.Equal(t, TaskDefinitionKind, instanceID.Kind)
			assert.Equal(t, tt.taskDefinition, instanceID.Name)
			assert.Equal(t, tt.containerName, instanceID.ContainerName)
			assert.Equal(t, container, instanceID.InstanceType)
			assert.Equal(t, tt.revision, instanceID.TaskDefinitionRevision)

			_, err := instanceID.GetSlug(false)
			require.NoError(t, err)

			expectedTemplateHash := fmt.Sprintf("%s:%s", tt.taskDefinition, tt.revision)
			expectedLabels := map[string]string{
				helpers.ApiGroupMetadataKey:                 "ecs",
				helpers.ApiVersionMetadataKey:               "2014-11-13",
				helpers.NamespaceMetadataKey:                tt.clusterName,
				helpers.KindMetadataKey:                     "TaskDefinition",
				helpers.NameMetadataKey:                     tt.taskDefinition,
				helpers.ContainerNameMetadataKey:            tt.containerName,
				helpers.TemplateHashKey:                     expectedTemplateHash,
				"kubescape.io/ecs-task-definition":          tt.taskDefinition,
				"kubescape.io/ecs-task-definition-revision": tt.revision,
			}

			labels := instanceID.GetLabels()
			assert.Equal(t, expectedLabels, labels)

			expectedFormatted := fmt.Sprintf("apiVersion-ecs/v1/namespace-%s/kind-TaskDefinition/name-%s/containerName-%s", tt.clusterName, tt.taskDefinition, tt.containerName)
			formatted := instanceID.GetStringFormatted()
			assert.Equal(t, expectedFormatted, formatted)

			expectedNoContainer := fmt.Sprintf("apiVersion-ecs/v1/namespace-%s/kind-TaskDefinition/name-%s", tt.clusterName, tt.taskDefinition)
			noContainer := instanceID.GetStringNoContainer()
			assert.Equal(t, expectedNoContainer, noContainer)
		})
	}
}

func TestECSInstanceID_Service(t *testing.T) {
	tests := []struct {
		name           string
		clusterName    string
		serviceName    string
		containerName  string
		taskDefinition string
		revision       string
	}{
		{
			name:           "basic service",
			clusterName:    "prod-cluster",
			serviceName:    "web-service",
			containerName:  "app",
			taskDefinition: "web-service",
			revision:       "5",
		},
		{
			name:           "service with different task definition name",
			clusterName:    "dev-cluster",
			serviceName:    "api-gateway",
			containerName:  "gateway",
			taskDefinition: "gateway-app",
			revision:       "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instanceID := NewServiceInstanceID(tt.clusterName, tt.serviceName, tt.containerName, tt.taskDefinition, tt.revision)

			assert.Equal(t, ECSDefaultApiVersion, instanceID.ApiVersion)
			assert.Equal(t, tt.clusterName, instanceID.ClusterName)
			assert.Equal(t, ServiceKind, instanceID.Kind)
			assert.Equal(t, tt.serviceName, instanceID.Name)
			assert.Equal(t, tt.containerName, instanceID.ContainerName)
			assert.Equal(t, container, instanceID.InstanceType)
			assert.Equal(t, tt.taskDefinition, instanceID.TaskDefinition)
			assert.Equal(t, tt.revision, instanceID.TaskDefinitionRevision)

			slug, err := instanceID.GetSlug(false)
			require.NoError(t, err)
			assert.NotEmpty(t, slug)

			expectedTemplateHash := fmt.Sprintf("%s:%s", tt.taskDefinition, tt.revision)
			expectedLabels := map[string]string{
				helpers.ApiGroupMetadataKey:                 "ecs",
				helpers.ApiVersionMetadataKey:               "2014-11-13",
				helpers.NamespaceMetadataKey:                tt.clusterName,
				helpers.KindMetadataKey:                     "Service",
				helpers.NameMetadataKey:                     tt.serviceName,
				helpers.ContainerNameMetadataKey:            tt.containerName,
				helpers.TemplateHashKey:                     expectedTemplateHash,
				"kubescape.io/ecs-task-definition":          tt.taskDefinition,
				"kubescape.io/ecs-task-definition-revision": tt.revision,
			}

			labels := instanceID.GetLabels()
			assert.Equal(t, expectedLabels, labels)
		})
	}
}

func TestECSInstanceID_Task(t *testing.T) {
	tests := []struct {
		name           string
		clusterName    string
		taskArn        string
		containerName  string
		taskDefinition string
		revision       string
	}{
		{
			name:           "basic task with ARN",
			clusterName:    "default",
			taskArn:        "arn:aws:ecs:us-east-1:123456789012:task/default/abcd1234-5678-90ef-ghij-klmnopqrstuv",
			containerName:  "worker",
			taskDefinition: "worker-task",
			revision:       "3",
		},
		{
			name:           "task with short ARN",
			clusterName:    "my-cluster",
			taskArn:        "arn:aws:ecs:us-west-2:987654321098:task/my-cluster/task-abc-123",
			containerName:  "processor",
			taskDefinition: "batch-job",
			revision:       "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instanceID := NewTaskInstanceID(tt.clusterName, tt.taskArn, tt.containerName, tt.taskDefinition, tt.revision)

			assert.Equal(t, ECSDefaultApiVersion, instanceID.ApiVersion)
			assert.Equal(t, tt.clusterName, instanceID.ClusterName)
			assert.Equal(t, TaskKind, instanceID.Kind)
			assert.Equal(t, tt.taskArn, instanceID.Name)
			assert.Equal(t, tt.containerName, instanceID.ContainerName)
			assert.Equal(t, container, instanceID.InstanceType)

			expectedAltName := extractTaskIDFromArn(tt.taskArn)
			assert.Equal(t, expectedAltName, instanceID.AlternateName)

			slug, err := instanceID.GetSlug(false)
			require.NoError(t, err)
			assert.NotEmpty(t, slug)

			expectedTemplateHash := fmt.Sprintf("%s:%s", tt.taskDefinition, tt.revision)
			expectedLabels := map[string]string{
				helpers.ApiGroupMetadataKey:                 "ecs",
				helpers.ApiVersionMetadataKey:               "2014-11-13",
				helpers.NamespaceMetadataKey:                tt.clusterName,
				helpers.KindMetadataKey:                     "Task",
				helpers.NameMetadataKey:                     tt.taskArn,
				helpers.ContainerNameMetadataKey:            tt.containerName,
				helpers.TemplateHashKey:                     expectedTemplateHash,
				"kubescape.io/ecs-task-definition":          tt.taskDefinition,
				"kubescape.io/ecs-task-definition-revision": tt.revision,
			}

			labels := instanceID.GetLabels()
			assert.Equal(t, expectedLabels, labels)
		})
	}
}

func TestECSInstanceID_GetOneTimeSlug(t *testing.T) {
	tests := []struct {
		name           string
		instanceID     *ECSInstanceID
		noContainer    bool
		prefixExpected string
	}{
		{
			name:           "task definition one-time slug",
			instanceID:     NewTaskDefinitionInstanceID("my-cluster", "app", "1", "container"),
			noContainer:    false,
			prefixExpected: "taskdefinition-app-container",
		},
		{
			name:           "service one-time slug",
			instanceID:     NewServiceInstanceID("prod", "api", "api-container", "api-task", "5"),
			noContainer:    false,
			prefixExpected: "service-api-api-container",
		},
		{
			name:           "task one-time slug with alternate name",
			instanceID:     NewTaskInstanceID("default", "arn:aws:ecs:us-east-1:123456789012:task/default/task-xyz-123", "worker", "worker-task", "2"),
			noContainer:    false,
			prefixExpected: "task-task-xyz-123-worker",
		},
		{
			name:           "task definition one-time slug no container",
			instanceID:     NewTaskDefinitionInstanceID("my-cluster", "app", "1", "container"),
			noContainer:    true,
			prefixExpected: "taskdefinition-app",
		},
		{
			name:           "service one-time slug no container",
			instanceID:     NewServiceInstanceID("prod", "api", "api-container", "api-task", "5"),
			noContainer:    true,
			prefixExpected: "service-api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oneTimeSlug, err := tt.instanceID.GetOneTimeSlug(tt.noContainer)
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(oneTimeSlug, tt.prefixExpected),
				"oneTimeSlug %q should start with %q", oneTimeSlug, tt.prefixExpected)

			assert.Greater(t, len(oneTimeSlug), len(tt.prefixExpected),
				"oneTimeSlug should be longer than prefix (includes UUID)")

			slug, err := tt.instanceID.GetSlug(tt.noContainer)
			require.NoError(t, err)

			lastHyphenIndex := strings.LastIndex(oneTimeSlug, "-")
			oneTimeSlugPrefix := oneTimeSlug[:lastHyphenIndex]

			assert.True(t, strings.HasPrefix(oneTimeSlugPrefix, strings.Split(slug, "-")[0]),
				"oneTimeSlug prefix should match slug prefix")
		})
	}
}

func TestECSInstanceID_GetSlug_NoContainer(t *testing.T) {
	tests := []struct {
		name         string
		instanceID   *ECSInstanceID
		noContainer  bool
		expectedSlug string
	}{
		{
			name:         "task definition without container",
			instanceID:   NewTaskDefinitionInstanceID("cluster", "task", "1", "app"),
			noContainer:  true,
			expectedSlug: "taskdefinition-task",
		},
		{
			name:         "service without container",
			instanceID:   NewServiceInstanceID("cluster", "svc", "web", "task-def", "3"),
			noContainer:  true,
			expectedSlug: "service-svc",
		},
		{
			name:         "task without container",
			instanceID:   NewTaskInstanceID("cluster", "arn:aws:ecs:us-east-1:123:task/cluster/task-id", "worker", "worker-task", "1"),
			noContainer:  true,
			expectedSlug: "task-task-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := tt.instanceID.GetSlug(tt.noContainer)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedSlug, slug)
		})
	}
}

func TestECSInstanceID_GetHashed(t *testing.T) {
	instanceID := NewTaskDefinitionInstanceID("my-cluster", "my-app", "1", "my-container")

	hashed := instanceID.GetHashed()

	assert.Equal(t, 64, len(hashed), "SHA256 hash should be 64 hex characters")
	assert.NotEmpty(t, hashed, "hashed value should not be empty")

	for _, c := range hashed {
		assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
			"hash should contain only hex characters, got: %c", c)
	}
}

func TestECSInstanceID_GetFormattedStrings(t *testing.T) {
	instanceID := NewServiceInstanceID("prod-cluster", "web-service", "nginx", "web-task", "7")

	formatted := instanceID.GetStringFormatted()
	expectedFormatted := "apiVersion-ecs/v1/namespace-prod-cluster/kind-Service/name-web-service/containerName-nginx"
	assert.Equal(t, expectedFormatted, formatted)

	noContainer := instanceID.GetStringNoContainer()
	expectedNoContainer := "apiVersion-ecs/v1/namespace-prod-cluster/kind-Service/name-web-service"
	assert.Equal(t, expectedNoContainer, noContainer)
}

func TestECSInstanceID_WithName(t *testing.T) {
	tests := []struct {
		name         string
		instanceID   *ECSInstanceID
		expectedName string
	}{
		{
			name:         "task definition returns family:revision",
			instanceID:   NewTaskDefinitionInstanceID("cluster", "my-app", "5", "container"),
			expectedName: "my-app",
		},
		{
			name:         "service returns service name",
			instanceID:   NewServiceInstanceID("cluster", "my-service", "container", "task-def", "1"),
			expectedName: "my-service",
		},
		{
			name:         "task returns ARN",
			instanceID:   NewTaskInstanceID("cluster", "arn:aws:ecs:us-east-1:123:task/cluster/task-123", "container", "task-def", "1"),
			expectedName: "arn:aws:ecs:us-east-1:123:task/cluster/task-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := tt.instanceID.GetName()
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

func TestECSInstanceID_WithAWSLabels(t *testing.T) {
	instanceID := NewServiceInstanceID("us-east-1-cluster", "api", "container", "api-task", "2")
	instanceID.Region = "us-east-1"
	instanceID.AccountID = "123456789012"

	labels := instanceID.GetLabels()

	assert.Equal(t, "us-east-1", labels["kubescape.io/aws-region"])
	assert.Equal(t, "123456789012", labels["kubescape.io/aws-account-id"])
	assert.Equal(t, "api-task", labels["kubescape.io/ecs-task-definition"])
	assert.Equal(t, "2", labels["kubescape.io/ecs-task-definition-revision"])
}

func TestExtractTaskIDFromArn(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "UUID task ID",
			arn:      "arn:aws:ecs:us-east-1:123456789012:task/default/abcd1234-5678-90ef-ghij-klmnopqrstuv",
			expected: "abcd1234-5678-90ef-ghij-klmnopqrstuv",
		},
		{
			name:     "Simple task ID",
			arn:      "arn:aws:ecs:us-west-2:987654321098:task/my-cluster/task-abc-123",
			expected: "task-abc-123",
		},
		{
			name:     "ARN without slash (edge case)",
			arn:      "task-xyz-789",
			expected: "task-xyz-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTaskIDFromArn(tt.arn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestECSInstanceID_TemplateHash(t *testing.T) {
	instanceID := NewTaskDefinitionInstanceID("cluster", "app", "10", "container")

	assert.Equal(t, "app:10", instanceID.GetTemplateHash())
	assert.Equal(t, "app:10", instanceID.TemplateHash)

	instanceID2 := NewServiceInstanceID("cluster", "service", "container", "task-def", "25")

	assert.Equal(t, "task-def:25", instanceID2.GetTemplateHash())
}

func TestECSInstanceID_InstanceType(t *testing.T) {
	instanceID := NewTaskDefinitionInstanceID("cluster", "app", "1", "container")

	assert.Equal(t, helpers.InstanceType(container), instanceID.GetInstanceType())
}

func TestECSInstanceID_GetContainerName(t *testing.T) {
	tests := []struct {
		name         string
		instanceID   *ECSInstanceID
		expectedName string
	}{
		{
			name:         "task definition container name",
			instanceID:   NewTaskDefinitionInstanceID("cluster", "app", "1", "nginx"),
			expectedName: "nginx",
		},
		{
			name:         "service container name",
			instanceID:   NewServiceInstanceID("cluster", "api", "gateway", "task", "1"),
			expectedName: "gateway",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, tt.instanceID.GetContainerName())
		})
	}
}

func TestECSInstanceID_Comprehensive(t *testing.T) {
	taskDefInstance := NewTaskDefinitionInstanceID("production", "payment-service", "15", "application")

	assert.Equal(t, "ecs/v1", taskDefInstance.ApiVersion)
	assert.Equal(t, "production", taskDefInstance.ClusterName)
	assert.Equal(t, "TaskDefinition", taskDefInstance.Kind)
	assert.Equal(t, "payment-service", taskDefInstance.Name)
	assert.Equal(t, "application", taskDefInstance.ContainerName)
	assert.Equal(t, "container", taskDefInstance.InstanceType)
	assert.Equal(t, "payment-service:15", taskDefInstance.TemplateHash)
	assert.Equal(t, "payment-service", taskDefInstance.TaskDefinition)
	assert.Equal(t, "15", taskDefInstance.TaskDefinitionRevision)

	slug, err := taskDefInstance.GetSlug(false)
	require.NoError(t, err)
	assert.NotEmpty(t, slug)
	assert.Contains(t, slug, "taskdefinition")
	assert.Contains(t, slug, "application")

	oneTimeSlug, err := taskDefInstance.GetOneTimeSlug(false)
	require.NoError(t, err)
	assert.NotEmpty(t, oneTimeSlug)
	assert.True(t, len(oneTimeSlug) > len(slug))
	assert.True(t, strings.HasPrefix(oneTimeSlug, strings.Split(slug, "-")[0]))

	labels := taskDefInstance.GetLabels()
	assert.Equal(t, "ecs", labels[helpers.ApiGroupMetadataKey])
	assert.Equal(t, "2014-11-13", labels[helpers.ApiVersionMetadataKey])
	assert.Equal(t, "production", labels[helpers.NamespaceMetadataKey])
	assert.Equal(t, "TaskDefinition", labels[helpers.KindMetadataKey])
	assert.Equal(t, "payment-service", labels[helpers.NameMetadataKey])
	assert.Equal(t, "application", labels[helpers.ContainerNameMetadataKey])
	assert.Equal(t, "payment-service:15", labels[helpers.TemplateHashKey])
}
