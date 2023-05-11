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
			"Empty image name returns empty value and error",
			"",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac8c",
			"",
			ErrInvalidFriendlyName,
		},
		{
			"Empty image hash returns empty value and error",
			"nginx",
			"",
			"",
			ErrInvalidFriendlyName,
		},
		{
			"Short image hash returns empty value and error",
			"nginx",
			"3ac8c",
			"",
			ErrInvalidFriendlyName,
		},
		{
			"Colon in image hash returns empty value and error",
			"nginx",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac:8c",
			"",
			ErrInvalidFriendlyName,
		},
		{
			"Slash in image hash returns empty value and error",
			"nginx",
			"f4e3b6489888647ce1834b601c6c06b9f8c03dee6e097e13ed3e28c01ea3ac/8c",
			"",
			ErrInvalidFriendlyName,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ImageInfoToFriendlyName(tc.imageTag, tc.imageHash)

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
			want:           "default-Pod-reverse-proxy-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "valid instanceID produces matching display name",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "webapp",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "default-Service-webapp-1ba5-4aaf",
			wantErr:        nil,
		},
		{
			name:           "invalid instanceID produces matching error",
			inputNamespace: "default",
			inputKind:      "Service",
			inputName:      "web/app",
			inputHashedID:  "1ba506b28f9ee9c7e8a0c98840fe5a1fe21142d225ecc526fbb535d0d6344aaf",
			want:           "",
			wantErr:        ErrInvalidFriendlyName,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := InstanceIDToFriendlyName(tc.inputName, tc.inputNamespace, tc.inputKind, tc.inputHashedID)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, tc.wantErr, err)
		})
	}
}
