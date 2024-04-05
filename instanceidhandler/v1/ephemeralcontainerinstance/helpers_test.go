package ephemeralcontainerinstance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	core1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
					apiVersion:             "",
					namespace:              "test",
					kind:                   "test",
					name:                   "test",
					ephemeralContainerName: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty namespace",
			args: args{
				instanceID: &InstanceID{
					apiVersion:             "test",
					namespace:              "",
					kind:                   "test",
					name:                   "test",
					ephemeralContainerName: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty kind",
			args: args{
				instanceID: &InstanceID{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "",
					name:                   "test",
					ephemeralContainerName: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty name",
			args: args{
				instanceID: &InstanceID{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "test",
					name:                   "",
					ephemeralContainerName: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "empty ephemeralContainerName",
			args: args{
				instanceID: &InstanceID{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "test",
					name:                   "test",
					ephemeralContainerName: "",
				},
			},
			wantErr: true,
		},
		{
			name: "valid instanceID",
			args: args{
				instanceID: &InstanceID{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "test",
					name:                   "test",
					ephemeralContainerName: "test",
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
		ownerReferences []metav1.OwnerReference
		containers      []core1.EphemeralContainer
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
				ownerReferences: []metav1.OwnerReference{},
				containers: []core1.EphemeralContainer{
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test",
						},
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "test",
			},
			want: []*InstanceID{
				{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "Pod",
					name:                   "test",
					ephemeralContainerName: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "empty containers",
			args: args{
				ownerReferences: []metav1.OwnerReference{},
				containers:      []core1.EphemeralContainer{},
				apiVersion:      "test",
				namespace:       "test",
				kind:            "test",
				name:            "test",
			},
			want:    []*InstanceID{},
			wantErr: false,
		},
		{
			name: "invalid instanceID",
			args: args{
				ownerReferences: []metav1.OwnerReference{},
				containers:      []core1.EphemeralContainer{},
				apiVersion:      "",
				namespace:       "test",
				kind:            "test",
				name:            "test",
			},
			want:    []*InstanceID{},
			wantErr: false,
		},
		{
			name: "valid instanceID - Pod",
			args: args{
				ownerReferences: []metav1.OwnerReference{},
				containers: []core1.EphemeralContainer{
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test",
						},
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "test",
			},
			want: []*InstanceID{
				{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "Pod",
					name:                   "test",
					ephemeralContainerName: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - Node",
			args: args{
				ownerReferences: []metav1.OwnerReference{
					{
						Kind: "Node",
						Name: "nodeName",
					},
				},
				containers: []core1.EphemeralContainer{
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test",
						},
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					apiVersion:             "test",
					namespace:              "test",
					kind:                   "Pod",
					name:                   "podName",
					ephemeralContainerName: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - multiple containers",
			args: args{
				ownerReferences: []metav1.OwnerReference{
					{
						APIVersion: "apps/v1",
						Kind:       "ReplicaSet",
						Name:       "OwnerTest",
					},
				},
				containers: []core1.EphemeralContainer{
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test",
						},
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					apiVersion:             "apps/v1",
					namespace:              "test",
					kind:                   "ReplicaSet",
					name:                   "OwnerTest",
					ephemeralContainerName: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "valid instanceID - Replica",
			args: args{
				ownerReferences: []metav1.OwnerReference{
					{
						Kind:       "ReplicaSet",
						Name:       "OwnerTest",
						APIVersion: "apps/v1",
					},
				},
				containers: []core1.EphemeralContainer{
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test-0",
						},
					},
					{
						EphemeralContainerCommon: core1.EphemeralContainerCommon{
							Name: "test-1",
						},
					},
				},
				apiVersion: "test",
				namespace:  "test",
				kind:       "Pod",
				name:       "podName",
			},
			want: []*InstanceID{
				{
					apiVersion:             "apps/v1",
					namespace:              "test",
					kind:                   "ReplicaSet",
					name:                   "OwnerTest",
					ephemeralContainerName: "test-0",
				},
				{
					apiVersion:             "apps/v1",
					namespace:              "test",
					kind:                   "ReplicaSet",
					name:                   "OwnerTest",
					ephemeralContainerName: "test-1",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listInstanceIDs(tt.args.ownerReferences, tt.args.containers, tt.args.apiVersion, tt.args.namespace, tt.args.kind, tt.args.name)
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
