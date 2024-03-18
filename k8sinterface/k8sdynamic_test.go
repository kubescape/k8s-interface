package k8sinterface

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubescape/k8s-interface/workloadinterface"
	//
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func TestListDynamic(t *testing.T) {
	if !IsConnectedToCluster() {
		return
	}
	k8s := NewKubernetesApi()
	// resource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	ww, err := k8s.ListWorkloads2("nginx-ingress", "Deployment")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ww) == 0 {
		return
	}
	s, _ := ww[0].GetSelector()
	g, _ := GetGroupVersionResource("pods")
	w, err := k8s.ListWorkloads(&g, "nginx-ingress", s.MatchLabels, nil)
	if len(w) != 1 {
		t.Error("empty list")
		return
	}

}

func mockWorkload(apiVersion, kind, namespace, name, ownerReferenceKind string) workloadinterface.IWorkload {
	mock := workloadinterface.NewWorkloadMock(nil)
	mock.SetKind(kind)
	mock.SetApiVersion(apiVersion)
	mock.SetName(name)
	mock.SetNamespace(namespace)

	if ownerReferenceKind != "" {
		apiVersion := ""
		switch ownerReferenceKind {
		case "Deployment", "ReplicaSet":
			apiVersion = "apps/v1"
		case "CronJob":
			apiVersion = "batch/v1"
		}
		ownerreferences := []metav1.OwnerReference{
			{
				APIVersion: apiVersion,
				Kind:       ownerReferenceKind,
			},
		}
		workloadinterface.SetInMap(mock.GetWorkload(), []string{"metadata"}, "ownerReferences", ownerreferences)
	}

	return mock
}
func TestWorkloadHasParent(t *testing.T) {

	// 1. Check if a nil workload returns false
	if WorkloadHasParent(nil) {
		t.Error("Expected false when provided a nil workload, but got true")
	}

	// 2. Provide an unsupported kind, ensure it returns false
	mockUnsupportedKind := mockWorkload("", "mockUnsupportedKind", "", "", "")
	if WorkloadHasParent(mockUnsupportedKind) {
		t.Error("Expected false for unsupported kind, but got true")
	}

	// 3. Provide a supported kind but no owner references
	mockJobWithoutOwner := mockWorkload("", "Job", "", "", "")
	if WorkloadHasParent(mockJobWithoutOwner) {
		t.Error("Expected false for Job without owner, but got true")
	}

	// 4. Provide a supported kind with owner references
	mockJobWithOwner := mockWorkload("", "Pod", "", "", "ReplicaSet")
	if !WorkloadHasParent(mockJobWithOwner) {
		t.Error("Expected true for Job with owner, but got false")
	}

	// 5. Provide a Pod without the pod-template-hash label
	mockPodWithoutHash := mockWorkload("", "Pod", "", "", "")
	mockPodWithoutHash.SetLabel("some-label", "value")
	if WorkloadHasParent(mockPodWithoutHash) {
		t.Error("Expected false for Pod without pod-template-hash, but got true")
	}

	// 6. Provide a Pod with the pod-template-hash label
	mockPodWithHash := mockWorkload("", "Pod", "", "", "")
	mockPodWithHash.SetLabel("pod-template-hash", "value")
	if !WorkloadHasParent(mockPodWithHash) {
		t.Error("Expected true for Pod with pod-template-hash, but got false")
	}

	// 7. Provide a Pod with an empty pod-template-hash label
	mockPodWithEmptyHash := mockWorkload("", "Pod", "", "", "")
	mockPodWithEmptyHash.SetLabel("pod-template-hash", "")
	if WorkloadHasParent(mockPodWithEmptyHash) {
		t.Error("Expected false for Pod with empty pod-template-hash, but got true")
	}

	// 8. Provide a ReplicaSet with owner references
	mockReplicaSetWithOwner := mockWorkload("", "ReplicaSet", "", "", "Deployment")
	if !WorkloadHasParent(mockReplicaSetWithOwner) {
		t.Error("Expected true for ReplicaSet with owner, but got false")
	}

	// 9. Provide a ReplicaSet without owner references
	mockReplicaSetWithoutOwner := mockWorkload("", "ReplicaSet", "", "", "")
	if WorkloadHasParent(mockReplicaSetWithoutOwner) {
		t.Error("Expected false for ReplicaSet without owner, but got true")
	}

	// 10. Provide a Job with multiple owner references
	mockJobWithMultipleOwners := mockWorkload("", "Job", "", "", "CronJob")
	if !WorkloadHasParent(mockJobWithMultipleOwners) {
		t.Error("Expected true for Job with multiple owners, but got false")
	}

	// 11. Test with other kinds of workloads to ensure they return false when they aren't "Pod", "Job", or "ReplicaSet".
	mockDaemonSet := mockWorkload("", "DaemonSet", "", "", "owner1")
	if WorkloadHasParent(mockDaemonSet) {
		t.Error("Expected false for DaemonSet regardless of owner references, but got true")
	}

	// 12. Test with a Pod that has both an owner reference and a pod-template-hash. It should return true because it satisfies both conditions.
	mockPodWithBoth := mockWorkload("", "Pod", "", "", "ReplicaSet")
	mockPodWithBoth.SetLabel("pod-template-hash", "value")
	if !WorkloadHasParent(mockPodWithBoth) {
		t.Error("Expected true for Pod with both owner reference and pod-template-hash, but got false")
	}

	// 13. Test a Pod without labels, it should return true.
	mockPodNoLabels := mockWorkload("", "Pod", "", "", "ReplicaSet")
	if !WorkloadHasParent(mockPodNoLabels) {
		t.Error("Expected false for Pod without labels, but got true")
	}

	// 14. Provide a ReplicaSet with owner references CRD
	mockReplicaSetWithCRDOwner := mockWorkload("", "ReplicaSet", "", "", "CRD")
	if WorkloadHasParent(mockReplicaSetWithCRDOwner) {
		t.Error("Expected true for ReplicaSet with owner, but got false")
	}

	// 15. Provide a Pod with owner references CRD
	mockPodWithCRDOwner := mockWorkload("", "Pod", "", "", "Node")
	if WorkloadHasParent(mockPodWithCRDOwner) {
		t.Error("Expected true for ReplicaSet with owner, but got false")
	}
}
