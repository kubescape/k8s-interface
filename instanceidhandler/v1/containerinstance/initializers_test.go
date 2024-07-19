package containerinstance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateInstanceIDFromString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    *InstanceID
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				input: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid input",
			args: args{
				input: "apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerMeme-nginx",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid input",
			args: args{
				input: "apiVersion-v1/namespace-default/kind-Pod/name-n/ginx/containerMeme-n/ginx",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid input - Pod",
			args: args{
				input: "apiVersion-v1/namespace-default/kind-Pod/name-nginx/containerName-nginx",
			},
			want: &InstanceID{
				ApiVersion:    "v1",
				Namespace:     "default",
				Kind:          "Pod",
				Name:          "nginx",
				ContainerName: "nginx",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - ReplicaSet",
			args: args{
				input: "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-1234/containerName-nginx",
			},
			want: &InstanceID{
				ApiVersion:    "apps/v1",
				Namespace:     "default",
				Kind:          "ReplicaSet",
				Name:          "nginx-1234",
				ContainerName: "nginx",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob",
			args: args{
				input: "apiVersion-batch/v1/namespace-kubescape/kind-Job/name-kubevuln-scheduler-b449cf78f/containerName-kubevuln-scheduler",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "kubescape",
				Kind:          "Job",
				Name:          "kubevuln-scheduler-b449cf78f", // should be kubevuln-scheduler-28677846
				ContainerName: "kubevuln-scheduler",
				InstanceType:  container,
				AlternateName: "", // should be kubevuln-scheduler-b449cf78f
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateInstanceIDFromString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateInstanceIDFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateInstanceIDFromString() = %v, want %v", got, tt.want)
		})
	}
}
