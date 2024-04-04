package ephemeralcontainerinstance

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/workloadinterface"
	"github.com/stretchr/testify/assert"
	core1 "k8s.io/api/core/v1"
)

// TestEphemeralGenerateInstanceIDFromString tests the instance id initialization
func TestEphemeralGenerateInstanceIDFromString(t *testing.T) {
	wp, err := workloadinterface.NewWorkload([]byte(mockPod))
	if err != nil {
		t.Fatalf(err.Error())
	}
	insFromWorkload, err := GenerateInstanceID(wp)
	if err != nil {
		t.Fatalf("can't create instance ID from pod")
	}

	p := &core1.Pod{}
	if err := json.Unmarshal([]byte(mockPod), p); err != nil {
		t.Fatalf(err.Error())
	}
	insFromPod, err := GenerateInstanceIDFromPod(p)
	if err != nil {
		t.Fatalf("can't create instance ID from pod")
	}

	assert.NotEqual(t, 0, len(insFromWorkload))
	assert.Equal(t, len(insFromWorkload), len(insFromPod))

	for i := range insFromWorkload {
		compare(t, &insFromWorkload[i], &insFromPod[i])
	}

	insFromString, err := GenerateInstanceIDFromString("apiVersion-v1/namespace-default/kind-Pod/name-nginx/ephemeralContainerName-nginx") //insFromWorkload[0].GetStringFormatted())
	if err != nil {
		t.Fatalf("can't create instance ID from string: %s, error: %s", insFromWorkload[0].GetStringFormatted(), err.Error())
	}
	compare(t, &insFromWorkload[0], insFromString)

}

func compare(t *testing.T, a, b instanceidhandler.IInstanceID) {
	assert.Equal(t, a.GetHashed(), b.GetHashed())
	assert.Equal(t, a.GetStringFormatted(), b.GetStringFormatted())

	assert.Equal(t, a.GetAPIVersion(), b.GetAPIVersion())
	assert.Equal(t, a.GetNamespace(), b.GetNamespace())
	assert.Equal(t, a.GetKind(), b.GetKind())
	assert.Equal(t, a.GetName(), b.GetName())
	assert.Equal(t, a.GetContainerName(), b.GetContainerName())
}

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
				input: "apiVersion-v1/namespace-default/kind-Pod/name-nginx/ephemeralContainerName-nginx",
			},
			want: &InstanceID{
				apiVersion:             "v1",
				namespace:              "default",
				kind:                   "Pod",
				name:                   "nginx",
				ephemeralContainerName: "nginx",
			},
			wantErr: false,
		},
		{
			name: "valid input - ReplicaSet",
			args: args{
				input: "apiVersion-apps/v1/namespace-default/kind-ReplicaSet/name-nginx-1234/ephemeralContainerName-nginx",
			},
			want: &InstanceID{
				apiVersion:             "apps/v1",
				namespace:              "default",
				kind:                   "ReplicaSet",
				name:                   "nginx-1234",
				ephemeralContainerName: "nginx",
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
			if got != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateInstanceIDFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
