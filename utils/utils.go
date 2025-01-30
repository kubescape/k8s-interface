package utils

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/scheme"
)

func ConvertUnstructuredToRuntimeObject(unstructuredObj *unstructured.Unstructured) (runtime.Object, error) {
	obj, err := scheme.Scheme.New(schema.FromAPIVersionAndKind(unstructuredObj.GetAPIVersion(), unstructuredObj.GetKind()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new object for apiVersion=%s kind=%s: %w", unstructuredObj.GetAPIVersion(), unstructuredObj.GetKind(), err)
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, obj); err != nil {
		return nil, fmt.Errorf("failed to convert unstructured to object: %w", err)
	}

	if _, ok := obj.(metav1.Object); !ok {
		return nil, fmt.Errorf("object is not a metav1.Object")
	}

	return obj, nil
}

func MustConvertRawToUnstructured(rawObject []byte) *unstructured.Unstructured {
	if x, err := ConvertRawToUnstructured(rawObject); err != nil {
		panic(err)
	} else {
		return x
	}
}

func MustConvertRuntimeObjectToUnstructured(obj runtime.Object) *unstructured.Unstructured {
	if x, err := ConvertRuntimeObjectToUnstructured(obj); err != nil {
		panic(err)
	} else {
		return x
	}
}

func ConvertRawToUnstructured(rawObject []byte) (*unstructured.Unstructured, error) {
	unstructuredObj := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawObject, unstructuredObj); err != nil {
		return nil, err
	}
	return unstructuredObj, nil
}

func ConvertRuntimeObjectToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: unstructuredObj}, nil
}
