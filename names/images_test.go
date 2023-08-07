package names

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeImageName(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "image tag",
			args: args{
				image: "nginx:latest",
			},
			want: "docker.io/library/nginx:latest",
		},
		{
			name: "image sha",
			args: args{
				image: "nginx@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
			},
			want: "docker.io/library/nginx@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
		},
		{
			name: "image tag sha",
			args: args{
				image: "nginx:latest@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
			},
			want: "docker.io/library/nginx:latest@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
		},
		{
			name: "repo image tag",
			args: args{
				image: "docker.io/library/nginx:latest",
			},
			want: "docker.io/library/nginx:latest",
		},
		{
			name: "repo image sha",
			args: args{
				image: "docker.io/library/nginx@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
			},
			want: "docker.io/library/nginx@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
		},
		{
			name: "repo image tag sha",
			args: args{
				image: "docker.io/library/nginx:latest@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
			},
			want: "docker.io/library/nginx:latest@sha256:73e957703f1266530db0aeac1fd6a3f87c1e59943f4c13eb340bb8521c6041d7",
		},
		{
			name: "quay image tag",
			args: args{
				image: "quay.io/kubescape/kubevuln:latest",
			},
			want: "quay.io/kubescape/kubevuln:latest",
		},
		{
			name: "quay image sha",
			args: args{
				image: "quay.io/kubescape/kubevuln@sha256:616d1d4312551b94088deb6ddab232ecabbbff0c289949a0d5f12d4b527c3f8a",
			},
			want: "quay.io/kubescape/kubevuln@sha256:616d1d4312551b94088deb6ddab232ecabbbff0c289949a0d5f12d4b527c3f8a",
		},
		{
			name: "quay image tag sha",
			args: args{
				image: "quay.io/kubescape/kubevuln:latest@sha256:616d1d4312551b94088deb6ddab232ecabbbff0c289949a0d5f12d4b527c3f8a",
			},
			want: "quay.io/kubescape/kubevuln:latest@sha256:616d1d4312551b94088deb6ddab232ecabbbff0c289949a0d5f12d4b527c3f8a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeImageName(tt.args.image)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
