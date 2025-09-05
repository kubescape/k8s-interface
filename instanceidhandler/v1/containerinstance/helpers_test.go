package containerinstance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const container = "container"

func Test_validateInstanceID(t *testing.T) {
	type args struct {
		instanceID *InstanceID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty instanceID",
			args: args{
				instanceID: &InstanceID{},
			},
			wantErr: true,
		},
		{
			name: "empty apiVersion",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "",
					Namespace:     "test",
					Kind:          "test",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: true,
		},
		{
			name: "empty namespace",
			args: args{
				instanceID: &InstanceID{

					ApiVersion:    "test",
					Namespace:     "",
					Kind:          "test",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: true,
		},
		{
			name: "empty kind",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: true,
		},
		{
			name: "empty name",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "test",
					Name:          "",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: true,
		},
		{
			name: "empty containerName",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "test",
					Name:          "test",
					ContainerName: "",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
		{
			name: "empty InstanceType",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "test",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  "",
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID",
			args: args{
				instanceID: &InstanceID{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "test",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInstanceID(tt.args.instanceID); (err != nil) != tt.wantErr {
				t.Errorf("validateInstanceID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_listInstanceIDs(t *testing.T) {
	type args struct {
		ownerReferences *metav1.OwnerReference
		containers      []corev1.Container
		apiVersion      string
		namespace       string
		kind            string
		name            string
	}
	tests := []struct {
		name    string
		args    args
		want    []*InstanceID
		wantErr bool
	}{
		{
			name: "empty ownerReferences",
			args: args{
				containers: []corev1.Container{
					{
						Name: "test",
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "test",
			},
			want: []*InstanceID{
				{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "Pod",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
		{
			name: "empty containers",
			args: args{
				containers: []corev1.Container{},
				apiVersion: "test",
				namespace:  "test",
				kind:       "test",
				name:       "test",
			},
			want:    []*InstanceID{},
			wantErr: false,
		},
		{
			name: "invalid instanceID",
			args: args{
				containers: []corev1.Container{
					{
						Name: "test",
					},
				},
				apiVersion: "",
				namespace:  "test",
				kind:       "test",
				name:       "test",
			},
			want:    []*InstanceID{},
			wantErr: true,
		},
		{
			name: "valid instanceID - Pod",
			args: args{
				containers: []corev1.Container{
					{
						Name: "test",
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "test",
			},
			want: []*InstanceID{
				{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "Pod",
					Name:          "test",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - Node",
			args: args{
				ownerReferences: &metav1.OwnerReference{
					Kind: "Node",
					Name: "nodeName",
				},
				containers: []corev1.Container{
					{
						Name: "test",
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					ApiVersion:    "test",
					Namespace:     "test",
					Kind:          "Pod",
					Name:          "podName",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - multiple containers",
			args: args{
				ownerReferences: &metav1.OwnerReference{
					APIVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Name:       "OwnerTest",
				},
				containers: []corev1.Container{
					{
						Name: "test",
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					ApiVersion:    "apps/v1",
					Namespace:     "test",
					Kind:          "ReplicaSet",
					Name:          "OwnerTest",
					ContainerName: "test",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - Replica",
			args: args{
				ownerReferences: &metav1.OwnerReference{
					Kind:       "ReplicaSet",
					Name:       "OwnerTest",
					APIVersion: "apps/v1",
				},
				containers: []corev1.Container{
					{
						Name: "test-0",
					},
					{
						Name: "test-1",
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					ApiVersion:    "apps/v1",
					Namespace:     "test",
					Kind:          "ReplicaSet",
					Name:          "OwnerTest",
					ContainerName: "test-0",
					InstanceType:  container,
				},
				{
					ApiVersion:    "apps/v1",
					Namespace:     "test",
					Kind:          "ReplicaSet",
					Name:          "OwnerTest",
					ContainerName: "test-1",
					InstanceType:  container,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListInstanceIDs(tt.args.ownerReferences, tt.args.containers, container, tt.args.apiVersion, tt.args.namespace, tt.args.kind, tt.args.name, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("listInstanceIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.want) != len(got) {
				t.Errorf("listInstanceIDs() len(tt.want) != len(got): %d != %d", len(tt.want), len(got))
				return
			}

			for i := range got {
				assert.Equal(t, tt.want[i].GetStringFormatted(), got[i].GetStringFormatted())
				assert.Equal(t, tt.want[i].GetHashed(), got[i].GetHashed())
			}
		})
	}
}
