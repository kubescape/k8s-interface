package apis

import (
	"testing"
)

func TestIsTypeDescribeRepositories(t *testing.T) {
	tests := []struct {
		name string
		obj  map[string]interface{}
		want bool
	}{
		{
			name: "valid DescribeRepositories",
			obj: map[string]interface{}{
				"apiVersion": "eks.amazonaws.com/v1",
				"kind":       "DescribeRepositories",
			},
			want: true,
		},
		{
			name: "invalid DescribeRepositories",
			obj: map[string]interface{}{
				"apiVersion": "eks.amazonaws.com/v1",
				"kind":       "Describe",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTypeDescribeRepositories(tt.obj); got != tt.want {
				t.Errorf("IsTypeDescribeRepositories(%v) = %v, want %v", tt.obj, got, tt.want)
			}
		})
	}
}

func TestIsTypeDescriptiveInfoFromCloudProvider(t *testing.T) {
	tests := []struct {
		name string
		obj  map[string]interface{}
		want bool
	}{
		{
			name: "valid ClusterDescribe",
			obj: map[string]interface{}{
				"apiVersion": "management.azure.com/v1",
				"kind":       "ClusterDescribe",
			},
			want: true,
		},
		{
			name: "invalid ClusterDescribe",
			obj: map[string]interface{}{
				"apiVersion": "container.googleapis.com/v1",
				"kind":       "InvalidType",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTypeDescriptiveInfoFromCloudProvider(tt.obj); got != tt.want {
				t.Errorf("IsTypeDescriptiveInfoFromCloudProvider(%v) = %v, want %v", tt.obj, got, tt.want)
			}
		})
	}
}

func TestIsType(t *testing.T) {
	tests := []struct {
		name            string
		object          map[string]interface{}
		acceptableTypes []string
		want            bool
	}{
		{
			name: "valid object and types",
			object: map[string]interface{}{
				"apiVersion": "management.azure.com/v1",
				"kind":       "ClusterDescribe",
			},
			acceptableTypes: []string{CloudProviderDescribeKind, "Describe"},
			want:            true,
		},
		{
			name: "invalid type",
			object: map[string]interface{}{
				"apiVersion": "container.googleapis.com/v1",
				"kind":       "InvalidType",
			},
			acceptableTypes: []string{CloudProviderDescribeKind, "Describe"},
			want:            false,
		},
		{
			name: "invalid apiVersion",
			object: map[string]interface{}{
				"apiVersion": "container/googleapis.com/v1",
				"kind":       "ClusterDescribe",
			},
			acceptableTypes: []string{CloudProviderDescribeKind, "Describe"},
			want:            false,
		},
		{
			name:            "nil object",
			object:          nil,
			acceptableTypes: []string{CloudProviderDescribeRepositoriesKind},
			want:            false,
		},
		{
			name: "valid DescribeRepositories",
			object: map[string]interface{}{
				"apiVersion": "eks.amazonaws.com/v1",
				"kind":       "DescribeRepositories",
			},
			acceptableTypes: []string{CloudProviderDescribeRepositoriesKind},
			want:            true,
		},
		{
			name: "valid ListEntitiesForPolicies",
			object: map[string]interface{}{
				"apiVersion": "eks.amazonaws.com/v1",
				"kind":       "ListEntitiesForPolicies",
			},
			acceptableTypes: []string{CloudProviderListEntitiesForPoliciesKind},
			want:            true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := IsCloudProviderType(tc.object, tc.acceptableTypes)
			if res != tc.want {
				t.Errorf("IsCloudProviderType(%v, %v) = %v, want %v", tc.object, tc.acceptableTypes, res, tc.want)
			}
		})
	}
}
