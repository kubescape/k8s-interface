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

func TestIsUnixTimeInMinutes(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"28677846", true},
		{"0", false},
		{"-1", false},
		{"b449cf78f", false},
		{"job", false},
	}
	for _, tt := range tests {
		if got := IsUnixTimeInMinutes(tt.input); got != tt.want {
			t.Errorf("IsUnixTimeInMinutes(%s) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsTemplateHash(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"b449cf78f", true},
		{"8698448884", true}, // entirely numeric hash generated from SafeEncodeString
		{"5678", true},       // short hash (< 8 chars)
		{"4", true},          // single character hash
		{"456789bcdf", true}, // full producer alphabet
		{"5f99858564", true},
		{"77bdd46fc5", true},
		{"84f5585d68", true},
		{"28677846", false},    // timestamp with digits '2' (not in producer domain)
		{"xxxxxxxx", false},    // non-domain character 'x'
		{"backup", false},      // non-domain word
		{"job", false},         // non-domain word
		{"nginx", false},       // non-domain word
		{"worker", false},      // non-domain word
		{"1234", false},        // digits '1','2','3' not in producer domain
		{"0456", false},        // digit '0' not in producer domain
		{"456789bcdf4", false}, // exceeds uint32 max length (11 chars)
		{"", false},            // empty string
	}
	for _, tt := range tests {
		if got := IsTemplateHash(tt.input); got != tt.want {
			t.Errorf("IsTemplateHash(%s) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
