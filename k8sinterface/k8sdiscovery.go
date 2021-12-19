package k8sinterface

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ValueNotFound = -1

var ResourceGroupMapping = map[string]string{}
var ResourceClusterScope = []string{}    // DEPRECATED
var ResourceNamesapcedScope = []string{} // use this to determan if the resource is namespaced

func InitializeMapResources(discoveryClient discovery.DiscoveryInterface) {

	// resourceList, _ := discoveryClient.ServerPreferredResources()
	// if len(resourceList) != 0 {
	// 	setMapResources(resourceList)
	// }

	// set mock initialization (if resources where missing from discovery. this can happen when an error accurse while pulling the resources)
	InitializeMapResourcesMock()

}
func setMapResources(resourceList []*metav1.APIResourceList) {
	for i := range resourceList {
		if resourceList[i] == nil {
			continue
		}
		if len(resourceList[i].APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(resourceList[i].GroupVersion)
		if err != nil {
			continue
		}
		if StringInSlice(ignoreGroups(), gv.Group) != ValueNotFound {
			continue
		}
		for _, apiResource := range resourceList[i].APIResources {
			if len(apiResource.Verbs) == 0 {
				continue
			}
			if _, ok := ResourceGroupMapping[apiResource.Name]; ok { // do not override resources in map
				continue
			}
			ResourceGroupMapping[apiResource.Name] = JoinGroupVersion(gv.Group, gv.Version)
			if apiResource.Namespaced {
				ResourceNamesapcedScope = append(ResourceNamesapcedScope, JoinResourceTriplets(gv.Group, gv.Version, apiResource.Name))
			} else { // DEPRECATED
				ResourceClusterScope = append(ResourceClusterScope, JoinResourceTriplets(gv.Group, gv.Version, apiResource.Name))

			}
		}
	}
}

func GetGroupVersionResource(resource string) (schema.GroupVersionResource, error) {
	resource = updateResourceKind(resource)
	if r, ok := ResourceGroupMapping[resource]; ok {
		gv := strings.Split(r, "/")
		return schema.GroupVersionResource{Group: gv[0], Version: gv[1], Resource: resource}, nil
	}
	return schema.GroupVersionResource{}, fmt.Errorf("resource '%s' unknown. Make sure the resource is found at `kubectl api-resources`", resource)
}

func IsNamespaceScope(resource *schema.GroupVersionResource) bool {
	return StringInSlice(ResourceNamesapcedScope, GroupVersionResourceToString(resource)) != ValueNotFound
}

func StringInSlice(strSlice []string, str string) int {
	for i := range strSlice {
		if strSlice[i] == str {
			return i
		}
	}
	return ValueNotFound
}

func JoinGroupVersion(group, version string) string {
	return fmt.Sprintf("%s/%s", group, version)
}

func JoinResourceTriplets(group, version, resource string) string {
	return fmt.Sprintf("%s/%s/%s", group, version, resource)
}

func GroupVersionResourceToString(resource *schema.GroupVersionResource) string {
	return JoinResourceTriplets(resource.Group, resource.Version, resource.Resource)
}
func GetResourceTriplets(group, version, resource string) []string {
	resourceTriplets := []string{}
	if resource == "" {
		// load full map
		for k, v := range ResourceGroupMapping {
			g := strings.Split(v, "/")
			resourceTriplets = append(resourceTriplets, JoinResourceTriplets(g[0], g[1], k))
		}
	} else if version == "" {
		// load by resource
		if v, ok := ResourceGroupMapping[resource]; ok {
			g := strings.Split(v, "/")
			if group == "" {
				group = g[0]
			}
			resourceTriplets = append(resourceTriplets, JoinResourceTriplets(group, g[1], resource))
		} else {
			// glog.Errorf("Resource '%s' unknown", resource)
		}
	} else if group == "" {
		// load by resource and version
		if v, ok := ResourceGroupMapping[resource]; ok {
			g := strings.Split(v, "/")
			resourceTriplets = append(resourceTriplets, JoinResourceTriplets(g[0], version, resource))
		} else {
			// glog.Errorf("Resource '%s' unknown", resource)
		}
	} else {
		resourceTriplets = append(resourceTriplets, JoinResourceTriplets(group, version, resource))
	}
	return resourceTriplets
}

// DEPRECATED
func ResourceGroupToString(group, version, resource string) []string {
	return ResourceGroupToSlice(group, version, resource)
}

func ResourceGroupToSlice(group, version, resource string) []string {
	if group == "*" {
		group = ""
	}
	if version == "*" {
		version = ""
	}
	if resource == "*" {
		resource = ""
	}
	resource = updateResourceKind(resource)
	return GetResourceTriplets(group, version, resource)
}

func StringToResourceGroup(str string) (string, string, string) {
	splitted := strings.Split(str, "/")
	for i := range splitted {
		if splitted[i] == "*" {
			splitted[i] = ""
		}
	}
	return splitted[0], splitted[1], splitted[2]
}

func updateResourceKind(resource string) string {
	resource = strings.ToLower(resource)

	if resource != "" && !strings.HasSuffix(resource, "s") {
		if strings.HasSuffix(resource, "y") {
			return fmt.Sprintf("%sies", strings.TrimSuffix(resource, "y")) // e.g. NetworkPolicy -> networkpolicies
		} else {
			return fmt.Sprintf("%ss", resource) // add 's' at the end of a resource
		}
	}
	return resource

}

func ignoreGroups() []string {
	return []string{"metrics.k8s.io"}
}
