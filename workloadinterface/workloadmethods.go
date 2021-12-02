package workloadinterface

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armosec/armoapi-go/apis"
	"github.com/armosec/utils-go/utils"
	"github.com/armosec/utils-k8s-go/armometadata"
	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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

// TODO - consider using a k8s manifest validator
// Return if this object is a valide k8s workload
func IsTypeWorkload(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if _, ok := object["apiVersion"]; !ok {
		return false
	}
	if _, ok := object["kind"]; !ok {
		return false
	}
	if imetadata, ok := object["metadata"]; ok {
		if metadata, ok := imetadata.(map[string]interface{}); ok {
			if _, ok := metadata["name"]; !ok {
				return false
			}
			// DO NOT TEST NAMESPACE - Not all k8s workloads have a namespace
		} else {
			return false
		}

	} else {
		return false
	}
	return true
}

// ======================================= DELETE ========================================

func (w *Workload) RemoveInject() {
	w.RemovePodLabel(armometadata.CAInject)      // DEPRECATED
	w.RemovePodLabel(armometadata.CAAttachLabel) // DEPRECATED
	w.RemovePodLabel(armometadata.ArmoAttach)

	w.RemoveLabel(armometadata.CAInject)      // DEPRECATED
	w.RemoveLabel(armometadata.CAAttachLabel) // DEPRECATED
	w.RemoveLabel(armometadata.ArmoAttach)
}

func (w *Workload) RemoveIgnore() {
	w.RemovePodLabel(armometadata.CAIgnore) // DEPRECATED
	w.RemovePodLabel(armometadata.ArmoAttach)

	w.RemoveLabel(armometadata.CAIgnore) // DEPRECATED
	w.RemoveLabel(armometadata.ArmoAttach)
}

func (w *Workload) RemoveWlid() {
	w.RemovePodAnnotation(armometadata.CAWlid) // DEPRECATED
	w.RemovePodAnnotation(armometadata.ArmoWlid)

	w.RemoveAnnotation(armometadata.CAWlid) // DEPRECATED
	w.RemoveAnnotation(armometadata.ArmoWlid)
}

func (w *Workload) RemoveCompatible() {
	w.RemovePodAnnotation(armometadata.ArmoCompatibleAnnotation)
}
func (w *Workload) RemoveJobID() {
	w.RemovePodAnnotation(armometadata.ArmoJobIDPath)
	w.RemovePodAnnotation(armometadata.ArmoJobParentPath)
	w.RemovePodAnnotation(armometadata.ArmoJobActionPath)

	w.RemoveAnnotation(armometadata.ArmoJobIDPath)
	w.RemoveAnnotation(armometadata.ArmoJobParentPath)
	w.RemoveAnnotation(armometadata.ArmoJobActionPath)
}
func (w *Workload) RemoveArmoMetadata() {
	w.RemoveArmoLabels()
	w.RemoveArmoAnnotations()
}

func (w *Workload) RemoveArmoAnnotations() {
	l := w.GetAnnotations()
	if l != nil {
		for k := range l {
			if strings.HasPrefix(k, armometadata.ArmoPrefix) {
				w.RemoveAnnotation(k)
			}
			if strings.HasPrefix(k, armometadata.CAPrefix) { // DEPRECATED
				w.RemoveAnnotation(k)
			}
		}
	}
	lp := w.GetPodAnnotations()
	if lp != nil {
		for k := range lp {
			if strings.HasPrefix(k, armometadata.ArmoPrefix) {
				w.RemovePodAnnotation(k)
			}
			if strings.HasPrefix(k, armometadata.CAPrefix) { // DEPRECATED
				w.RemovePodAnnotation(k)
			}
		}
	}
}
func (w *Workload) RemoveArmoLabels() {
	l := w.GetLabels()
	if l != nil {
		for k := range l {
			if strings.HasPrefix(k, armometadata.ArmoPrefix) {
				w.RemoveLabel(k)
			}
			if strings.HasPrefix(k, armometadata.CAPrefix) { // DEPRECATED
				w.RemoveLabel(k)
			}
		}
	}
	lp := w.GetPodLabels()
	if lp != nil {
		for k := range lp {
			if strings.HasPrefix(k, armometadata.ArmoPrefix) {
				w.RemovePodLabel(k)
			}
			if strings.HasPrefix(k, armometadata.CAPrefix) { // DEPRECATED
				w.RemovePodLabel(k)
			}
		}
	}
}
func (w *Workload) RemoveUpdateTime() {

	// remove from pod
	w.RemovePodAnnotation(armometadata.CAUpdate) // DEPRECATED
	w.RemovePodAnnotation(armometadata.ArmoUpdate)

	// remove from workload
	w.RemoveAnnotation(armometadata.CAUpdate) // DEPRECATED
	w.RemoveAnnotation(armometadata.ArmoUpdate)
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

func (w *Workload) SetKind(kind string) {
	w.workload["kind"] = kind
}

func (w *Workload) SetInject() {
	w.SetPodLabel(armometadata.ArmoAttach, utils.BoolToString(true))
}

func (w *Workload) SetJobID(jobTracking apis.JobTracking) {
	w.SetPodAnnotation(armometadata.ArmoJobIDPath, jobTracking.JobID)
	w.SetPodAnnotation(armometadata.ArmoJobParentPath, jobTracking.ParentID)
	w.SetPodAnnotation(armometadata.ArmoJobActionPath, fmt.Sprintf("%d", jobTracking.LastActionNumber))
}

func (w *Workload) SetIgnore() {
	w.SetPodLabel(armometadata.ArmoAttach, utils.BoolToString(false))
}

func (w *Workload) SetCompatible() {
	w.SetPodAnnotation(armometadata.ArmoCompatibleAnnotation, utils.BoolToString(true))
}

func (w *Workload) SetIncompatible() {
	w.SetPodAnnotation(armometadata.ArmoCompatibleAnnotation, utils.BoolToString(false))
}

func (w *Workload) SetReplaceheaders() {
	w.SetPodAnnotation(armometadata.ArmoReplaceheaders, utils.BoolToString(true))
}

func (w *Workload) SetWlid(wlid string) {
	w.SetPodAnnotation(armometadata.ArmoWlid, wlid)
}

func (w *Workload) SetUpdateTime() {
	w.SetPodAnnotation(armometadata.ArmoUpdate, string(time.Now().UTC().Format("02-01-2006 15:04:05")))
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
	return fmt.Sprintf("%s/%s/%s/%s", w.GetApiVersion(), w.GetNamespace(), w.GetKind(), w.GetName())
}
func (w *Workload) GetName() string {
	if v, ok := InspectWorkload(w.workload, "metadata", "name"); ok {
		return v.(string)
	}
	return ""
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
		replicas, isok := v.(float64)
		if isok {
			return int(replicas)
		}
		if replicas, isok := v.(int64); isok {
			return int(replicas)
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
	if v, ok := InspectWorkload(w.workload, "spec", "selector", "matchLabels"); ok && v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return selector, err
		}
		if err := json.Unmarshal(b, selector); err != nil {
			return selector, err
		}
		return selector, nil
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

// func (w *Workload) GetJobID() string {
// 	if status, ok := w.GetAnnotation(armometadata.ArmoJobID); ok {
// 		return status
// 	}
// 	return ""
// }

// ========================================= IS =========================================

func (w *Workload) IsInject() bool {
	return w.IsAttached()
}

func (w *Workload) IsIgnore() bool {
	if attach := armometadata.IsAttached(w.GetPodLabels()); attach != nil {
		return !(*attach)
	}
	if attach := armometadata.IsAttached(w.GetLabels()); attach != nil {
		return !(*attach)
	}
	return false
}

func (w *Workload) IsCompatible() bool {
	if c, ok := w.GetPodAnnotation(armometadata.ArmoCompatibleAnnotation); ok {
		return utils.StringToBool(c)

	}
	if c, ok := w.GetAnnotation(armometadata.ArmoCompatibleAnnotation); ok {
		return utils.StringToBool(c)

	}
	return false
}

func (w *Workload) IsIncompatible() bool {
	if c, ok := w.GetPodAnnotation(armometadata.ArmoCompatibleAnnotation); ok {
		return !utils.StringToBool(c)
	}
	if c, ok := w.GetAnnotation(armometadata.ArmoCompatibleAnnotation); ok {
		return !utils.StringToBool(c)
	}
	return false
}
func (w *Workload) IsAttached() bool {
	if attach := armometadata.IsAttached(w.GetPodLabels()); attach != nil {
		return *attach
	}
	if attach := armometadata.IsAttached(w.GetLabels()); attach != nil {
		return *attach
	}
	return false
}

func (w *Workload) IsReplaceheaders() bool {
	if c, ok := w.GetPodAnnotation(armometadata.ArmoReplaceheaders); ok {
		return utils.StringToBool(c)
	}
	return false
}
