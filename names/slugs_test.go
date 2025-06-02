package names

import (
	"testing"

	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
)

func TestImageInfoToFriendlyName(t *testing.T) {
	tt := []struct {
		name      string
		imageTag  string
		imageHash string
		expected  string
		wantErr   error
	}{
		{
			"Short image tag returns matching value",
			"nginx",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"nginx-a3ac8c",
			nil,
		},
		{
			"Short versioned image tag returns matching value",
			"nginx:latest",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"nginx-latest-a3ac8c",
			nil,
		},
		{
			"Full image tag returns matching value",
			"docker.io/nginx:latest",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"docker.io-nginx-latest-a3ac8c",
			nil,
		},
		{
			"Image ID format produces matching value",
			"docker-pullable://gcr.io/etcd-development/etcd",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"docker-pullable-gcr.io-etcd-development-etcd-a3ac8c",
			nil,
		},
		{
			"Image ID format with uppercase symbols produces matching lowercase value",
			"docker-pullable://GCR.io/etcD-development/Etcd",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"docker-pullable-gcr.io-etcd-development-etcd-a3ac8c",
			nil,
		},
		{
			"Image ID format with underscore works",
			"quay.io/matthiasb_1/kubevuln:renaming",
			"quay.io/matthiasb_1/kubevuln@sha256:85c1b06d541d61ddb46efcd8b316855f544278c9ab27a07ec35bbe81be54fbec",
			"quay.io-matthiasb-1-kubevuln-renaming-54fbec",
			nil,
		},
		{
			"Image ID format with at sign works",
			"docker.io/kindest/local-path-provisioner:v0.0.23-kind.0@sha256:f2d0a02831ff3a03cf51343226670d5060623b43a4cfc4808bd0875b2c4b9501",
			"docker.io/kindest/local-path-provisioner:v0.0.23-kind.0@sha256:f2d0a02831ff3a03cf51343226670d5060623b43a4cfc4808bd0875b2c4b9501",
			"docker.io-kindest-local-path-provisioner-v0.0.23-kind.0-sha256-f2d0a02831ff3a03cf51343226670d5060623b43a4cfc4808bd0875b2c4b9501-4b9501",
			nil,
		},
		{
			"Empty image name returns empty value and error",
			"",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"",
			ErrInvalidSlug,
		},
		{
			"Empty image hash returns empty value and error",
			"nginx",
			"",
			"",
			ErrInvalidSlug,
		},
		{
			"Short image hash returns empty value and error",
			"nginx",
			"3ac8c",
			"",
			ErrInvalidSlug,
		},
		{
			"Colon in image hash returns empty value and error",
			"nginx",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac:8c",
			"",
			ErrInvalidSlug,
		},
		{
			"Slash in image hash returns empty value and error",
			"nginx",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac/8c",
			"",
			ErrInvalidSlug,
		},
		{
			"Image names that would produce overflowing slugs should get truncated to limit",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabc",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab-a3ac8c",
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ImageInfoToSlug(tc.imageTag, tc.imageHash)

			assert.Equal(t, tc.expected, got)
			assert.ErrorIs(t, tc.wantErr, err)
		})
	}
}

func TestInstanceIDToFriendlyName(t *testing.T) {
	tt := []struct {
		name           string
		inputName      string
		inputNamespace string
		inputContainer string
		inputKind      string
		inputHashedID  string
		want           string
		wantErr        error
	}{
		{
			name:           "valid instanceID produces matching display name",
			inputNamespace: "default",
			inputKind:      "Pod",
			inputName:      "reverse-proxy",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "pod-reverse-proxy",
			wantErr:        nil,
		},
		{
			name:           "valid instanceID produces matching display name",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "webapp",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "service-webapp",
			wantErr:        nil,
		},
		{
			name:           "valid instanceID with container name produces matching display name",
			inputNamespace: "default",
			inputContainer: "webapp",
			inputKind:      "Service",
			inputName:      "webapp",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "service-webapp-webapp-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "valid instanceID with different namespace name produces different hash",
			inputNamespace: "blabla",
			inputContainer: "webapp",
			inputKind:      "Service",
			inputName:      "webapp",
			inputHashedID:  "000006b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6340000",
			want:           "service-webapp-webapp-0000-0000",
			wantErr:        nil,
		},
		{
			name:           "instanceID that produces overflowing slugs gets truncated to limit",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "service-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "invalid instanceID produces matching error",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "web/app",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "",
			wantErr:        ErrInvalidSlug,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := InstanceIDToSlug(tc.inputName, tc.inputKind, tc.inputContainer, tc.inputHashedID)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestIsValidSubdomainName(t *testing.T) {
	tt := []struct {
		name      string
		inputName string
		want      bool
	}{
		{
			name:      "Short alphanumeric name is considered valid",
			inputName: "nginx",
			want:      true,
		},
		{
			name:      "Colon character should render the string invalid",
			inputName: "n:ginx",
			want:      false,
		},
		{
			name:      "Slash character should render the string invalid",
			inputName: "n/ginx",
			want:      false,
		},
		{
			name:      "Caret character should render the string invalid",
			inputName: "n^ginx",
			want:      false,
		},
		{
			name:      "Empty string should be considered invalid",
			inputName: "",
			want:      false,
		},
		{
			name:      "Periods should be allowed",
			inputName: "docker.io",
			want:      true,
		},
		{
			name:      "Hyphens should be allowed",
			inputName: "web-app",
			want:      true,
		},
		{
			name:      "Numbers should be allowed",
			inputName: "webapp1",
			want:      true,
		},
		{
			name:      "Names starting from an allowed non-alphanumeric character should be invalid",
			inputName: "-webapp",
			want:      false,
		},
		{
			name:      "Names ending with an allowed non-alphanumeric character should be invalid",
			inputName: "webapp-",
			want:      false,
		},
		{
			name:      "Names over 253 characters should be invalid",
			inputName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaX",
			want:      false,
		},
		{
			name:      "An uppercase character is considered invalid",
			inputName: "nGinx",
			want:      false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidDNSSubdomainName(tc.inputName)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsValidDSNLabelName(t *testing.T) {
	tt := []struct {
		name      string
		inputName string
		want      bool
	}{
		{
			name:      "Short alphanumeric name is considered valid",
			inputName: "nginx",
			want:      true,
		},
		{
			name:      "Colon character should render the string invalid",
			inputName: "n:ginx",
			want:      false,
		},
		{
			name:      "Slash character should render the string invalid",
			inputName: "n/ginx",
			want:      false,
		},
		{
			name:      "Caret character should render the string invalid",
			inputName: "n^ginx",
			want:      false,
		},
		{
			name:      "Empty string should be considered invalid",
			inputName: "",
			want:      false,
		},
		{
			name:      "Periods should NOT be allowed",
			inputName: "docker.io",
			want:      false,
		},
		{
			name:      "Hyphens should be allowed",
			inputName: "web-app",
			want:      true,
		},
		{
			name:      "Numbers should be allowed",
			inputName: "webapp1",
			want:      true,
		},
		{
			name:      "Names starting from an allowed non-alphanumeric character should be invalid",
			inputName: "-webapp",
			want:      false,
		},
		{
			name:      "Names ending with an allowed non-alphanumeric character should be invalid",
			inputName: "webapp-",
			want:      false,
		},
		{
			name:      "Names over 63 characters should be invalid",
			inputName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaX",
			want:      false,
		},
		{
			name:      "An uppercase character is considered valid",
			inputName: "nginX",
			want:      false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidDNSLabelName(tc.inputName)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsValidLabelValue(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"", true},
		{"valid_value", true},
		{"Valid_value", true},
		{"Valid_value1", true},
		{"1Valid_value", true},
		{"valid.value", true},
		{"valid-value", true},
		{"valid_value", true},
		{"invalid:value", false},
		{"$special_char", false},
		{"very_long_value_that_is_more_than_63_characters_long_and_should_fail_validation", false},
	}

	for _, test := range tests {
		actual := IsValidLabelValue(test.value)
		if actual != test.expected {
			t.Errorf("For value '%s', expected %t but got %t", test.value, test.expected, actual)
		}
	}
}

func TestToValidLabelValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"valid_value", "valid_value"},
		{"1Valid.value", "1Valid.value"},
		{"1number", "1number"},
		{"$special:char", "special-char"},
		{"-very_long_value_that:is_more@than_63_characters_long_and_should_fail_validation", "very_long_value_that-is_more-than_63_characters_long_and_should"},
	}

	for _, test := range tests {
		actual := ToValidLabelValue(test.input)
		if actual != test.expected {
			t.Errorf("For input '%s', expected '%s' but got '%s'", test.input, test.expected, actual)
		}
	}
}

func TestToValidDNSSubdomainName(t *testing.T) {
	tt := []struct {
		name          string
		inputName     string
		want          string
		expectedError bool
	}{
		{
			name:      "Short alphanumeric name is considered valid",
			inputName: "nginx",
			want:      "nginx",
		},
		{
			name:      "Colon character should render the string invalid, and should be replaced with a hyphen",
			inputName: "n:ginx",
			want:      "n-ginx",
		},
		{
			name:      "Slash character should render the string invalid, and should be replaced with a hyphen",
			inputName: "n/ginx",
			want:      "n-ginx",
		},
		{
			name:      "Caret character should render the string invalid, and should be replaced with a hyphen",
			inputName: "n^ginx",
			want:      "n-ginx",
		},
		{
			name:          "Empty string should be considered invalid and should return an error",
			inputName:     "",
			want:          "",
			expectedError: true,
		},
		{
			name:      "Periods should be allowed",
			inputName: "docker.io",
			want:      "docker.io",
		},
		{
			name:      "Hyphens should be allowed",
			inputName: "web-app",
			want:      "web-app",
		},
		{
			name:      "Numbers should be allowed",
			inputName: "webapp1",
			want:      "webapp1",
		},
		{
			name:      "Names starting from an allowed non-alphanumeric character should be invalid and should be truncated",
			inputName: "-webapp",
			want:      "webapp",
		},
		{
			name:      "Names ending with an allowed non-alphanumeric character should be invalid and should be truncated",
			inputName: "webapp-",
			want:      "webapp",
		},
		{
			name:      "Names over 253 characters should be invalid and should be truncated to limit",
			inputName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaX",
			want:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		{
			name:      "An uppercase character is considered invalid and should be converted to lowercase",
			inputName: "nGinx",
			want:      "nginx",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToValidDNSSubdomainName(tc.inputName)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestStringToSlug(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "short input",
			input:    "n:ginx-xyz1.2.34",
			expected: "n-ginx-xyz1.2.34",
			err:      nil,
		},
		{
			name:     "long input",
			input:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			expected: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123-60e5-4c7d",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slug, err := StringToSlug(tc.input)
			if err != tc.err {
				t.Errorf("Expected error: %v, got: %v", tc.err, err)
			}
			if slug != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, slug)
			}
			assert.LessOrEqual(t, len(slug), MaxDNSSubdomainLength)
		})
	}
}

type FakeMetadata struct {
	workloadinterface.IMetadata

	Namespace  string
	ApiVersion string
	Kind       string
	Name       string
	ID         string
}

func (f *FakeMetadata) GetID() string {
	return f.ID
}

func (f *FakeMetadata) GetNamespace() string {
	return f.Namespace
}

func (f *FakeMetadata) GetApiVersion() string {
	return f.ApiVersion
}

func (f *FakeMetadata) GetKind() string {
	return f.Kind
}

func (f *FakeMetadata) GetName() string {
	return f.Name
}

func TestResourceToSlug(t *testing.T) {
	testCases := []struct {
		name     string
		resource workloadinterface.IMetadata
		expected string
	}{
		{
			name: "",
			resource: &FakeMetadata{
				ApiVersion: "v1",
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "mypod",
			},
			expected: "pod-mypod",
		},
		{
			resource: &FakeMetadata{
				ApiVersion: "",
				Kind:       "Pod",
				Namespace:  "",
				Name:       "mypod",
			},
			expected: "pod-mypod",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ResourceToSlug(tc.resource)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
			assert.NoError(t, err)
		})
	}
}

func TestRoleBindingResourceToSlug(t *testing.T) {
	testCases := []struct {
		name        string
		subject     workloadinterface.IMetadata
		role        workloadinterface.IMetadata
		roleBinding workloadinterface.IMetadata
		expected    string
	}{
		{
			name: "role, rolebinding",
			subject: &FakeMetadata{
				Kind:      "ServiceAccount",
				Name:      "sa-2",
				Namespace: "kubescape",
			},
			role: &FakeMetadata{
				Kind:      "Role",
				Name:      "myrole",
				Namespace: "namespace-1",
			},
			roleBinding: &FakeMetadata{
				Kind:      "RoleBinding",
				Name:      "myrolebinding",
				Namespace: "namespace-2",
			},
			expected: "serviceaccount-sa-2-role-myrole-rolebinding-myrolebinding",
		},
		{
			name: "with related objects (cluster role, cluster rolebinding)",
			subject: &FakeMetadata{
				Kind:      "ServiceAccount",
				Name:      "sa-1",
				Namespace: "kubescape",
			},
			role: &FakeMetadata{
				Kind: "ClusterRole",
				Name: "myrole",
			},
			roleBinding: &FakeMetadata{
				Kind: "ClusterRoleBinding",
				Name: "myrolebinding",
			},
			expected: "serviceaccount-sa-1-clusterrole-myrole-clusterrolebinding-myrolebinding",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := RoleBindingResourceToSlug(tc.subject, tc.role, tc.roleBinding)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
			assert.NoError(t, err)
		})
	}
}
