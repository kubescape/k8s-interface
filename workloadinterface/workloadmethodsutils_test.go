package workloadinterface

import (
	"fmt"
	"testing"
)

func TestGetReplicasetNameFromPod(t *testing.T) {
	tt := []struct {
		podName         string
		podTemplateHash string
		want            string
		wantErr         error
	}{
		{
			podName:         "pod-name-123-456",
			podTemplateHash: "123",
			want:            "pod-name-123",
			wantErr:         nil,
		},
		{
			podName:         "pod-name-789-012",
			podTemplateHash: "789",
			want:            "pod-name-789",
			wantErr:         nil,
		},
		{
			podName:         "bla-pod-name-abc2236-def",
			podTemplateHash: "abc2236",
			want:            "bla-pod-name-abc2236",
			wantErr:         nil,
		},
		{
			podName:         "pod-name",
			podTemplateHash: "123",
			want:            "",
			wantErr:         fmt.Errorf("failed to get replicaset name from pod name: pod-name"),
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("podName=%s, podTemplateHash=%s", tc.podName, tc.podTemplateHash), func(t *testing.T) {
			got, err := getReplicasetNameFromPod(tc.podName, tc.podTemplateHash)

			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}

			if (err != nil && tc.wantErr == nil) || (err == nil && tc.wantErr != nil) || (err != nil && tc.wantErr != nil && err.Error() != tc.wantErr.Error()) {
				t.Errorf("got error %v, want error %v", err, tc.wantErr)
			}
		})
	}
}
