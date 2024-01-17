package helpers

import "testing"

func Test_ignoreOwnerReference(t *testing.T) {
	type args struct {
		ownerKind string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ignore - Node",
			args: args{
				ownerKind: "Node",
			},
			want: true,
		},
		{
			name: "ignore - CRD",
			args: args{
				ownerKind: "Bla",
			},
			want: true,
		},
		{
			name: "not ignore - Pod",
			args: args{
				ownerKind: "Pod",
			},
			want: false,
		},
		{
			name: "not ignore",
			args: args{
				ownerKind: "ReplicaSet",
			},
			want: false,
		},
		{
			name: "not ignore - StatefulSet",
			args: args{
				ownerKind: "StatefulSet",
			},
			want: false,
		},
		{
			name: "not ignore - Job",
			args: args{
				ownerKind: "Job",
			},
			want: false,
		},
		{
			name: "not ignore - CronJob",
			args: args{
				ownerKind: "CronJob",
			},
			want: false,
		},
		{
			name: "not ignore - DaemonSet",
			args: args{
				ownerKind: "DaemonSet",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IgnoreOwnerReference(tt.args.ownerKind); got != tt.want {
				t.Errorf("ignoreOwnerReference() = %v, want %v", got, tt.want)
			}
		})
	}
}
