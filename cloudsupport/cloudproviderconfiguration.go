package cloudsupport

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	container "cloud.google.com/go/container/apiv1"
	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/armosec/k8s-interface/workloadinterface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

type CloudProviderDescription struct {
	object map[string]interface{}
}

const TypeCloudProviderDescription workloadinterface.ObjectType = "cloudProviderDescription"

const (
	CloudProviderGroup           = "cloudvendordata.armo.cloud"
	CloudProviderDescriptionKind = "description"
)

// Setters
func (obj *CloudProviderDescription) SetNamespace(namespace string) {
	obj.SetProvider(namespace)
}

func (obj *CloudProviderDescription) SetGroup(group string) {
	obj.object["group"] = group
}

func (obj *CloudProviderDescription) SetName(name string) {
	obj.object["name"] = name
}

func (obj *CloudProviderDescription) SetProvider(provider string) {
	obj.object["provider"] = provider
}

func (obj *CloudProviderDescription) SetKind(kind string) {
	obj.object["kind"] = kind
}

func (obj *CloudProviderDescription) SetWorkload(object map[string]interface{}) {
	obj.SetObject(object)
}

func (obj *CloudProviderDescription) SetObject(object map[string]interface{}) {
	obj.object = object
}

// Getters

//group -> cloudvendordata.armo.cloud
func (obj *CloudProviderDescription) GetGroup() string {
	if v, ok := workloadinterface.InspectMap(obj.object, "group"); ok {
		return v.(string)
	}
	return ""
}

func (obj *CloudProviderDescription) GetApiVersion() string {
	return fmt.Sprintf("%s/%s", obj.GetGroup(), "v1beta0")
}

func (obj *CloudProviderDescription) GetObjectType() workloadinterface.ObjectType {
	return TypeCloudProviderDescription
}
func (obj *CloudProviderDescription) GetKind() string {
	if v, ok := workloadinterface.InspectMap(obj.object, "kind"); ok {
		return v.(string)
	}
	return ""
}

func (obj *CloudProviderDescription) GetName() string {
	if v, ok := workloadinterface.InspectMap(obj.object, "name"); ok {
		return v.(string)
	}
	return ""
}

// provider -> eks/gke
func (obj *CloudProviderDescription) GetProvider() string {
	if v, ok := workloadinterface.InspectMap(obj.object, "provider"); ok {
		return v.(string)
	}
	return ""
}

func (obj *CloudProviderDescription) GetNamespace() string {
	return obj.GetProvider()
}

func (obj *CloudProviderDescription) GetWorkload() map[string]interface{} {
	return obj.GetObject()
}

func (obj *CloudProviderDescription) GetObject() map[string]interface{} {
	return obj.object
}

// cloudvendordata.armo.cloud/provider/description/clusterName
func (obj *CloudProviderDescription) GetID() string {
	return fmt.Sprintf("%s/%s/%s/%s", obj.GetGroup(), obj.GetProvider(), obj.GetKind(), obj.GetName())
}

func NewDescriptiveInfoFromCloudProvider(object map[string]interface{}) *CloudProviderDescription {
	return &CloudProviderDescription{
		object: object,
	}
}

func IsTypeDescriptiveInfoFromCloudProvider(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	if kind, ok := object["kind"]; !ok || kind != CloudProviderDescriptionKind {
		return false
	} else if _, ok := object["group"]; !ok {
		return false
	} else {
		if object["kind"] != CloudProviderDescriptionKind || object["group"] != CloudProviderGroup {
			return false
		}
	}

	if _, ok := object["name"]; !ok {
		return false
	}
	if _, ok := object["provider"]; !ok {
		return false
	}

	return true
}

func IsRunningInCloudProvider() bool {
	currContext := k8sinterface.GetCurrentContext()
	if currContext == nil {
		return false
	}
	if strings.Contains(currContext.Cluster, strings.ToLower("eks")) || strings.Contains(currContext.Cluster, strings.ToLower("gke")) {
		return true
	}
	return false
}

func GetCloudProvider(currContext string) string {
	if strings.Contains(currContext, strings.ToLower("eks")) {
		return "eks"
	} else if strings.Contains(currContext, strings.ToLower("gke")) {
		return "gke"
	}
	return ""
}

func GetDescriptiveInfoFromCloudProvider() (workloadinterface.IMetadata, error) {
	currContext := k8sinterface.GetCurrentContext()
	var clusterInfo *CloudProviderDescription
	var err error
	if currContext == nil {
		return nil, nil
	}
	cloudProvider := GetCloudProvider(currContext.Cluster)
	switch cloudProvider {
	case "eks":
		clusterInfo, err = GetClusterInfoForEKS(currContext)
	case "gke":
		clusterInfo, err = GetClusterInfoForGKE()
	}

	if err != nil {
		return nil, err
	}
	clusterInfo.SetKind(CloudProviderDescriptionKind)
	clusterInfo.SetGroup(CloudProviderGroup)
	return clusterInfo, nil
}

// Get descriptive info about cluster running in EKS.
func GetClusterInfoForEKS(currContext *api.Context) (*CloudProviderDescription, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	splittedClusterContext := strings.Split(k8sinterface.GetCurrentContext().Cluster, ".")
	if len(splittedClusterContext) < 2 {
		return nil, fmt.Errorf("error: failed to get region")
	}
	region := splittedClusterContext[1]
	// Configure cluster name and region for request
	svc := eks.New(s, &aws.Config{Region: aws.String(region)})
	input := &eks.DescribeClusterInput{
		Name: aws.String(k8sinterface.GetClusterName()),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		return nil, err
	}
	resultInJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	var clusterInfo *CloudProviderDescription
	err = json.Unmarshal(resultInJson, &clusterInfo.object)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetName(*result.Cluster.Name)
	clusterInfo.SetProvider("eks")
	return clusterInfo, nil
}

// Get descriptive info about cluster running in GKE.
func GetClusterInfoForGKE() (*CloudProviderDescription, error) {
	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	parsedName := strings.Split(k8sinterface.GetClusterName(), "_")
	if len(parsedName) < 3 {
		return nil, fmt.Errorf("error: failed to parse cluster name")
	}
	clusterName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", parsedName[1], parsedName[2], parsedName[3])
	req := &containerpb.GetClusterRequest{
		Name: clusterName,
	}
	result, err := c.GetCluster(ctx, req)
	if err != nil {
		return nil, err
	}
	resultInJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	var clusterInfo *CloudProviderDescription
	err = json.Unmarshal(resultInJson, &clusterInfo.object)
	if err != nil {
		return nil, err
	}
	clusterInfo.SetProvider("gke")
	return clusterInfo, nil
}
