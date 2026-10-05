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
			name: "valid input - CronJob child Job",
			args: args{
				input: "apiVersion-batch/v1/namespace-kubescape/kind-Job/name-kubevuln-scheduler-b449cf78f/containerName-kubevuln-scheduler",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "kubescape",
				Kind:          "Job",
				Name:          "kubevuln-scheduler-b449cf78f",
				ContainerName: "kubevuln-scheduler",
				InstanceType:  container,
				AlternateName: "kubevuln-scheduler-b449cf78f",
				TemplateHash:  "b449cf78f",
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob child Job, no container",
			args: args{
				input: "apiVersion-batch/v1/namespace-kubescape/kind-Job/name-kubevuln-scheduler-b449cf78f",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "kubescape",
				Kind:          "Job",
				Name:          "kubevuln-scheduler-b449cf78f",
				AlternateName: "kubevuln-scheduler-b449cf78f",
				TemplateHash:  "b449cf78f",
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob workload",
			args: args{
				input: "apiVersion-batch/v1/namespace-kubescape/kind-CronJob/name-kubevuln-scheduler/containerName-kubevuln-scheduler",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "kubescape",
				Kind:          "CronJob",
				Name:          "kubevuln-scheduler",
				ContainerName: "kubevuln-scheduler",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob workload, no container",
			args: args{
				input: "apiVersion-batch/v1/namespace-kubescape/kind-CronJob/name-kubevuln-scheduler",
			},
			want: &InstanceID{
				ApiVersion: "batch/v1",
				Namespace:  "kubescape",
				Kind:       "CronJob",
				Name:       "kubevuln-scheduler",
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob child Job with numeric template hash",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-8698448884/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-8698448884",
				ContainerName: "backup",
				InstanceType:  container,
				AlternateName: "backup-8698448884",
				TemplateHash:  "8698448884",
			},
			wantErr: false,
		},
		{
			name: "valid input - CronJob child Job with short template hash",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-5678/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-5678",
				ContainerName: "backup",
				InstanceType:  container,
				AlternateName: "backup-5678",
				TemplateHash:  "5678",
			},
			wantErr: false,
		},
		{
			name: "valid input - ordinary CronJob named backup-xxxxxxxx preserves literal name",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-CronJob/name-backup-xxxxxxxx/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "CronJob",
				Name:          "backup-xxxxxxxx",
				ContainerName: "backup",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - ordinary CronJob named backup-8698448884 preserves literal name",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-CronJob/name-backup-8698448884/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "CronJob",
				Name:          "backup-8698448884",
				ContainerName: "backup",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - standard Job with timestamp suffix preserves literal name",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-28677846/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-28677846",
				ContainerName: "backup",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - standard Job named backup-xxxxxxxx preserves literal name",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-xxxxxxxx/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-xxxxxxxx",
				ContainerName: "backup",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - literal Job backup-db parsed as template hash acknowledging tradeoff",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-db/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-db",
				ContainerName: "backup",
				InstanceType:  container,
				AlternateName: "backup-db",
				TemplateHash:  "db",
			},
			wantErr: false,
		},
		{
			name: "valid input - literal Job backup-4 parsed as template hash acknowledging tradeoff",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-backup-4/containerName-backup",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "backup-4",
				ContainerName: "backup",
				InstanceType:  container,
				AlternateName: "backup-4",
				TemplateHash:  "4",
			},
			wantErr: false,
		},
		{
			name: "valid input - standard Job without hash",
			args: args{
				input: "apiVersion-batch/v1/namespace-default/kind-Job/name-nginx-job/containerName-nginx-job",
			},
			want: &InstanceID{
				ApiVersion:    "batch/v1",
				Namespace:     "default",
				Kind:          "Job",
				Name:          "nginx-job",
				ContainerName: "nginx-job",
				InstanceType:  container,
			},
			wantErr: false,
		},
		{
			name: "valid input - ReplicaSet, no container",
			args: args{
				input: "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-1234",
			},
			want: &InstanceID{
				ApiVersion: "apps/v1",
				Namespace:  "default",
				Kind:       "ReplicaSet",
				Name:       "nginx-1234",
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
