package names

import (
	"testing"

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
			want:           "default-pod-reverse-proxy-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "valid instanceID produces matching display name",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "webapp",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "default-service-webapp-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "instanceID that produces overflowing slugs gets truncated to limit",
			inputNamespace: "0123456789",
			inputKind:      "0123456789",
			inputName:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "0123456789-0123456789-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-1ba5-4aaf",
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
			got, err := InstanceIDToSlug(tc.inputName, tc.inputNamespace, tc.inputKind, tc.inputHashedID)

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
