package cloudsupport

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	container "cloud.google.com/go/container/apiv1"
	k8sinterface "github.com/armosec/k8s-interface/k8sinterface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetDescriptiveInfoFromCloudProvider() (map[string]interface{}, error) {
	currContext := k8sinterface.GetCurrentContext()
	if currContext == nil {
		return nil, nil
	}
	if strings.Contains(currContext.Cluster, strings.ToLower("eks")) {
		return GetClusterInfoForEKS(currContext)
	} else if strings.Contains(currContext.Cluster, strings.ToLower("gke")) {
		return GetClusterInfoForGKE()
	}
	return nil, nil
}

// Get descriptive info about cluster running in EKS.
func GetClusterInfoForEKS(currContext *api.Context) (map[string]interface{}, error) {
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
	var clusterInfo map[string]interface{}
	err = json.Unmarshal(resultInJson, &clusterInfo)
	if err != nil {
		return nil, err
	}
	return clusterInfo, nil
}

// Get descriptive info about cluster running in GKE.
func GetClusterInfoForGKE() (map[string]interface{}, error) {
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
	var clusterInfo map[string]interface{}
	err = json.Unmarshal(resultInJson, &clusterInfo)
	if err != nil {
		return nil, err
	}
	return clusterInfo, nil
}
