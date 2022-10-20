package workloadinterface

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/armosec/armoapi-go/apis"
	"github.com/armosec/utils-k8s-go/armometadata"
	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/strings/slices"
)

const TypeWorkloadObject ObjectType = "workload"

type Workload struct {
	workload map[string]interface{}
}

func NewWorkload(bWorkload []byte) (*Workload, error) {
	workload := make(map[string]interface{})
	if bWorkload != nil {
		if err := json.Unmarshal(bWorkload, &workload); err != nil {
			return nil, err
		}
	}
	// if !IsTypeWorkload(workload) {
	// 	return nil, fmt.Errorf("invalid workload - expected k8s workload")
	// }
	return &Workload{
		workload: workload,
	}, nil
}

func NewWorkloadObj(workload map[string]interface{}) *Workload {
	return &Workload{
		workload: workload,
	}
}

func (w *Workload) GetObjectType() ObjectType {
	return TypeWorkloadObject
}

func (w *Workload) Json() string {
	return w.ToString()
}
func (w *Workload) ToString() string {
	if w.GetWorkload() == nil {
		return ""
	}
	bWorkload, err := json.Marshal(w.GetWorkload())
	if err != nil {
		return err.Error()
	}
	return string(bWorkload)
}

func (workload *Workload) DeepCopy(w map[string]interface{}) {
	workload.workload = make(map[string]interface{})
	byt, _ := json.Marshal(w)
	json.Unmarshal(byt, &workload.workload)
}

func (w *Workload) ToUnstructured() (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	if w.workload == nil {
		return obj, nil
	}
	bWorkload, err := json.Marshal(w.workload)
	if err != nil {
		return obj, err
	}
	if err := json.Unmarshal(bWorkload, obj); err != nil {
		return obj, err

	}

	return obj, nil
}

// ======================================= DELETE ========================================

func (w *Workload) RemoveJobID() {
	w.RemovePodAnnotation(armometadata.ArmoJobIDPath)
	w.RemovePodAnnotation(armometadata.ArmoJobParentPath)
	w.RemovePodAnnotation(armometadata.ArmoJobActionPath)

	w.RemoveAnnotation(armometadata.ArmoJobIDPath)
	w.RemoveAnnotation(armometadata.ArmoJobParentPath)
	w.RemoveAnnotation(armometadata.ArmoJobActionPath)
}

func (w *Workload) RemoveSecretData() {
	w.RemoveAnnotation("kubectl.kubernetes.io/last-applied-configuration")
	delete(w.workload, "data")
}

func (w *Workload) RemovePodStatus() {
	delete(w.workload, "status")
}

func (w *Workload) RemoveResourceVersion() {
	if _, ok := w.workload["metadata"]; !ok {
		return
	}
	meta, _ := w.workload["metadata"].(map[string]interface{})
	delete(meta, "resourceVersion")
}

func (w *Workload) RemoveLabel(key string) {
	w.RemoveMetadata([]string{"metadata"}, "labels", key)
}

func (w *Workload) RemoveAnnotation(key string) {
	w.RemoveMetadata([]string{"metadata"}, "annotations", key)
}

func (w *Workload) RemovePodAnnotation(key string) {
	w.RemoveMetadata(PodMetadata(w.GetKind()), "annotations", key)
}

func (w *Workload) RemovePodLabel(key string) {
	w.RemoveMetadata(PodMetadata(w.GetKind()), "labels", key)
}

func (w *Workload) RemoveMetadata(scope []string, metadata, key string) {

	workload := w.workload
	for i := range scope {
		if _, ok := workload[scope[i]]; !ok {
			return
		}
		workload, _ = workload[scope[i]].(map[string]interface{})
	}

	if _, ok := workload[metadata]; !ok {
		return
	}

	labels, _ := workload[metadata].(map[string]interface{})
	delete(labels, key)

}

// ========================================= SET =========================================

func (w *Workload) SetWorkload(workload map[string]interface{}) {
	w.SetObject(workload)
}

func (w *Workload) SetObject(workload map[string]interface{}) {
	w.workload = workload
}

func (w *Workload) SetApiVersion(apiVersion string) {
	w.workload["apiVersion"] = apiVersion
}

func (w *Workload) SetKind(kind string) {
	w.workload["kind"] = kind
}

func (w *Workload) SetJobID(jobTracking apis.JobTracking) {
	w.SetPodAnnotation(armometadata.ArmoJobIDPath, jobTracking.JobID)
	w.SetPodAnnotation(armometadata.ArmoJobParentPath, jobTracking.ParentID)
	w.SetPodAnnotation(armometadata.ArmoJobActionPath, fmt.Sprintf("%d", jobTracking.LastActionNumber))
}

func (w *Workload) SetNamespace(namespace string) {
	SetInMap(w.workload, []string{"metadata"}, "namespace", namespace)
}

func (w *Workload) SetName(name string) {
	SetInMap(w.workload, []string{"metadata"}, "name", name)
}

func (w *Workload) SetLabel(key, value string) {
	SetInMap(w.workload, []string{"metadata", "labels"}, key, value)
}

func (w *Workload) SetPodLabel(key, value string) {
	SetInMap(w.workload, append(PodMetadata(w.GetKind()), "labels"), key, value)
}
func (w *Workload) SetAnnotation(key, value string) {
	SetInMap(w.workload, []string{"metadata", "annotations"}, key, value)
}
func (w *Workload) SetPodAnnotation(key, value string) {
	SetInMap(w.workload, append(PodMetadata(w.GetKind()), "annotations"), key, value)
}

// ========================================= GET =========================================
func (w *Workload) GetWorkload() map[string]interface{} {
	return w.GetObject()
}
func (w *Workload) GetObject() map[string]interface{} {
	return w.workload
}
func (w *Workload) GetNamespace() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "namespace"); ok {
		return v.(string)
	}
	return ""
}
func (w *Workload) GetID() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", w.GetGroup(), w.GetVersion(), w.GetNamespace(), w.GetKind(), w.GetName())

	// TODO - return like selfLink - e.g. /apis/apps/v1/namespaces/monitoring/statefulsets/alertmanager-prometheus-
	// return fmt.Sprintf("apps/%s/%s/%s/%s", w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName())
}
func (w *Workload) GetName() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "name"); ok {
		return v.(string)
	}
	return ""
}

func (w *Workload) GetData() map[string]interface{} {
	if v, ok := InspectWorkload(w.workload, "data"); ok {
		return v.(map[string]interface{})
	}
	return nil
}

func (w *Workload) GetApiVersion() string {
	if v, ok := InspectWorkload(w.workload, "apiVersion"); ok {
		return v.(string)
	}
	return ""
}

func (w *Workload) GetVersion() string {
	apiVersion := w.GetApiVersion()
	splitted := strings.Split(apiVersion, "/")
	if len(splitted) == 1 {
		return splitted[0]
	} else if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func (w *Workload) GetGroup() string {
	apiVersion := w.GetApiVersion()
	splitted := strings.Split(apiVersion, "/")
	if len(splitted) == 2 {
		return splitted[0]
	}
	return ""
}

func (w *Workload) GetGenerateName() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "generateName"); ok {
		return v.(string)
	}
	return ""
}

func (w *Workload) GetReplicas() int {
	if v, ok := InspectWorkload(w.workload, "spec", "replicas"); ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int64:
			return int(n)
		case float32:
			return int(n)
		case int32:
			return int(n)
		case int16:
			return int(n)
		case int:
			return n
		}
	}
	return 1
}

func (w *Workload) GetKind() string {
	if v, ok := InspectWorkload(w.workload, "kind"); ok {
		return v.(string)
	}
	return ""
}
func (w *Workload) GetSelector() (*metav1.LabelSelector, error) {
	selector := &metav1.LabelSelector{}
	if matchLabels, ok := InspectWorkload(w.workload, "spec", "selector", "matchLabels"); ok && matchLabels != nil {
		if m, ok := matchLabels.(map[string]interface{}); ok {
			selector.MatchLabels = make(map[string]string, len(m))
			for k, v := range m {
				selector.MatchLabels[k] = v.(string)
			}
		}
	}
	if matchExpressions, ok := InspectWorkload(w.workload, "spec", "selector", "matchExpressions"); ok && matchExpressions != nil {
		b, err := json.Marshal(matchExpressions)
		if err != nil {
			return selector, err
		}
		if err := json.Unmarshal(b, &selector.MatchExpressions); err != nil {
			return selector, nil
		}
	}
	return selector, nil
}

func (w *Workload) GetAnnotation(annotation string) (string, bool) {
	if v, ok := InspectWorkload(w.workload, "metadata", "annotations", annotation); ok {
		return v.(string), ok
	}
	return "", false
}
func (w *Workload) GetLabel(label string) (string, bool) {
	if v, ok := InspectWorkload(w.workload, "metadata", "labels", label); ok {
		return v.(string), ok
	}
	return "", false
}

func (w *Workload) GetPodLabel(label string) (string, bool) {
	if v, ok := InspectWorkload(w.workload, append(PodMetadata(w.GetKind()), "labels", label)...); ok && v != nil {
		return v.(string), ok
	}
	return "", false
}

func (w *Workload) GetLabels() map[string]string {
	if v, ok := InspectWorkload(w.workload, "metadata", "labels"); ok && v != nil {
		labels := make(map[string]string)
		for k, i := range v.(map[string]interface{}) {
			labels[k] = i.(string)
		}
		return labels
	}
	return nil
}

// GetInnerLabels - DEPRECATED
func (w *Workload) GetInnerLabels() map[string]string {
	return w.GetPodLabels()
}

func (w *Workload) GetPodLabels() map[string]string {
	if v, ok := InspectWorkload(w.workload, append(PodMetadata(w.GetKind()), "labels")...); ok && v != nil {
		labels := make(map[string]string)
		for k, i := range v.(map[string]interface{}) {
			labels[k] = i.(string)
		}
		return labels
	}
	return nil
}

// GetInnerAnnotations - DEPRECATED
func (w *Workload) GetInnerAnnotations() map[string]string {
	return w.GetPodAnnotations()
}

// GetPodAnnotations
func (w *Workload) GetPodAnnotations() map[string]string {
	if v, ok := InspectWorkload(w.workload, append(PodMetadata(w.GetKind()), "annotations")...); ok && v != nil {
		annotations := make(map[string]string)
		for k, i := range v.(map[string]interface{}) {
			annotations[k] = fmt.Sprintf("%v", i)
		}
		return annotations
	}
	return nil
}

// GetInnerAnnotation DEPRECATED
func (w *Workload) GetInnerAnnotation(annotation string) (string, bool) {
	return w.GetPodAnnotation(annotation)
}

func (w *Workload) GetPodAnnotation(annotation string) (string, bool) {
	if v, ok := InspectWorkload(w.workload, append(PodMetadata(w.GetKind()), "annotations", annotation)...); ok && v != nil {
		return v.(string), ok
	}
	return "", false
}

func (w *Workload) GetAnnotations() map[string]string {
	if v, ok := InspectWorkload(w.workload, "metadata", "annotations"); ok && v != nil {
		annotations := make(map[string]string)
		for k, i := range v.(map[string]interface{}) {
			annotations[k] = fmt.Sprintf("%v", i)
		}
		return annotations
	}
	return nil
}

// GetVolumes -
func (w *Workload) GetVolumes() ([]corev1.Volume, error) {
	volumes := []corev1.Volume{}

	interVolumes, _ := InspectWorkload(w.workload, append(PodSpec(w.GetKind()), "volumes")...)
	if interVolumes == nil {
		return volumes, nil
	}
	volumesBytes, err := json.Marshal(interVolumes)
	if err != nil {
		return volumes, err
	}
	err = json.Unmarshal(volumesBytes, &volumes)

	return volumes, err
}

func (w *Workload) GetServiceAccountName() string {

	if v, ok := InspectWorkload(w.workload, append(PodSpec(w.GetKind()), "serviceAccountName")...); ok && v != nil {
		return v.(string)
	}
	return ""
}

func (w *Workload) GetPodSpec() (*corev1.PodSpec, error) {
	podSpec := &corev1.PodSpec{}
	podSepcRaw, _ := InspectWorkload(w.workload, PodSpec(w.GetKind())...)
	if podSepcRaw == nil {
		return podSpec, fmt.Errorf("no PodSpec for workload: %v", w)
	}
	b, err := json.Marshal(podSepcRaw)
	if err != nil {
		return podSpec, err
	}
	err = json.Unmarshal(b, podSpec)

	return podSpec, err
}

func (w *Workload) GetImagePullSecret() ([]corev1.LocalObjectReference, error) {
	imgPullSecrets := []corev1.LocalObjectReference{}

	iImgPullSecrets, _ := InspectWorkload(w.workload, append(PodSpec(w.GetKind()), "imagePullSecrets")...)
	b, err := json.Marshal(iImgPullSecrets)
	if err != nil {
		return imgPullSecrets, err
	}
	err = json.Unmarshal(b, &imgPullSecrets)

	return imgPullSecrets, err
}

// GetContainers -
func (w *Workload) GetContainers() ([]corev1.Container, error) {
	containers := []corev1.Container{}

	interContainers, _ := InspectWorkload(w.workload, append(PodSpec(w.GetKind()), "containers")...)
	if interContainers == nil {
		return containers, nil
	}
	containersBytes, err := json.Marshal(interContainers)
	if err != nil {
		return containers, err
	}
	err = json.Unmarshal(containersBytes, &containers)

	return containers, err
}

// GetInitContainers -
func (w *Workload) GetInitContainers() ([]corev1.Container, error) {
	containers := []corev1.Container{}

	interContainers, _ := InspectWorkload(w.workload, append(PodSpec(w.GetKind()), "initContainers")...)
	if interContainers == nil {
		return containers, nil
	}
	containersBytes, err := json.Marshal(interContainers)
	if err != nil {
		return containers, err
	}
	err = json.Unmarshal(containersBytes, &containers)

	return containers, err
}

// GetOwnerReferences -
func (w *Workload) GetOwnerReferences() ([]metav1.OwnerReference, error) {
	ownerReferences := []metav1.OwnerReference{}
	interOwnerReferences, ok := InspectWorkload(w.workload, "metadata", "ownerReferences")
	if !ok {
		return ownerReferences, nil
	}

	ownerReferencesBytes, err := json.Marshal(interOwnerReferences)
	if err != nil {
		return ownerReferences, err
	}
	err = json.Unmarshal(ownerReferencesBytes, &ownerReferences)
	if err != nil {
		return ownerReferences, err

	}
	return ownerReferences, nil
}
func (w *Workload) GetResourceVersion() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "resourceVersion"); ok {
		return v.(string)
	}
	return ""
}
func (w *Workload) GetUID() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "uid"); ok {
		return v.(string)
	}
	return ""
}
func (w *Workload) GetWlid() string {
	if wlid, ok := w.GetAnnotation(armometadata.ArmoWlid); ok {
		return wlid
	}
	return ""
}

func (w *Workload) GenerateWlid(clusterName string) string {
	return wlidpkg.GetK8sWLID(clusterName, w.GetNamespace(), w.GetKind(), w.GetName())
}

func (w *Workload) GetJobID() *apis.JobTracking {
	jobTracking := apis.JobTracking{}
	if job, ok := w.GetPodAnnotation(armometadata.ArmoJobIDPath); ok {
		jobTracking.JobID = job
	}
	if parent, ok := w.GetPodAnnotation(armometadata.ArmoJobParentPath); ok {
		jobTracking.ParentID = parent
	}
	if action, ok := w.GetPodAnnotation(armometadata.ArmoJobActionPath); ok {
		if i, err := strconv.Atoi(action); err == nil {
			jobTracking.LastActionNumber = i
		}
	}
	if jobTracking.LastActionNumber == 0 { // start the counter at 1
		jobTracking.LastActionNumber = 1
	}
	return &jobTracking
}

func (w *Workload) GetSecretsToContainer() (map[string][]string, error) {
	secretsToContainer := make(map[string][]string)
	workloadSecrets, err := w.GetSecrets()
	if err != nil {
		return nil, err
	}
	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		for _, volumeMount := range container.VolumeMounts {
			if slices.Contains(workloadSecrets, volumeMount.Name) {
				secretsToContainer[container.Name] = append(secretsToContainer[volumeMount.Name], container.Name)
			}
		}

	}
	return secretsToContainer, nil
}

func (w *Workload) GetConfigMapsToContainer() (map[string][]string, error) {
	configMapsToContainer := make(map[string][]string)
	workloadConfigMaps, err := w.GetConfigMaps()
	if err != nil {
		return nil, err
	}
	containers, err := w.GetContainers()
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		for _, volumeMount := range container.VolumeMounts {
			if slices.Contains(workloadConfigMaps, volumeMount.Name) {
				configMapsToContainer[container.Name] = append(configMapsToContainer[volumeMount.Name], container.Name)
			}
		}

	}
	return configMapsToContainer, nil
}

func (w *Workload) GetSecrets() ([]string, error) {
	volumes, err := w.GetVolumes()
	if err != nil {
		return nil, err
	}
	secrets := []string{}
	for _, volume := range volumes {
		if volume.Secret != nil {
			secrets = append(secrets, volume.Name)
		}
	}

	return secrets, nil
}

func (w *Workload) GetConfigMaps() ([]string, error) {
	volumes, err := w.GetVolumes()
	if err != nil {
		return nil, err
	}
	configMaps := []string{}
	for _, volume := range volumes {
		if volume.ConfigMap != nil {
			configMaps = append(configMaps, volume.Name)
		}
	}

	return configMaps, nil
}
