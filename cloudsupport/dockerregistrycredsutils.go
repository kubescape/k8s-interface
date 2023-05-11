package cloudsupport

import (
	"context"
	"fmt"

	logger "github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"

	"github.com/armosec/utils-k8s-go/secrethandling"
	"github.com/docker/docker/api/types"
	"github.com/kubescape/k8s-interface/k8sinterface"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listPodImagePullSecrets(podSpec *corev1.PodSpec) ([]string, error) {
	if podSpec == nil {
		return []string{}, fmt.Errorf("in listPodImagePullSecrets podSpec is nil")
	}
	secrets := []string{}
	for _, i := range podSpec.ImagePullSecrets {
		secrets = append(secrets, i.Name)
	}
	return secrets, nil
}

func listServiceAccountImagePullSecrets(k8sAPI *k8sinterface.KubernetesApi, namespace, serviceAccountName string) ([]string, error) {
	secrets := []string{}
	if serviceAccountName == "" {
		return secrets, nil
	}

	serviceAccount, err := k8sAPI.KubernetesClient.CoreV1().ServiceAccounts(namespace).Get(k8sAPI.Context, serviceAccountName, metav1.GetOptions{})
	if err != nil {
		return secrets, fmt.Errorf("in listServiceAccountImagePullSecrets failed to get ServiceAccounts: %v", err)
	}
	for i := range serviceAccount.ImagePullSecrets {
		secrets = append(secrets, serviceAccount.ImagePullSecrets[i].Name)
	}
	return secrets, nil
}

func getImagePullSecret(k8sAPI *k8sinterface.KubernetesApi, secrets []string, namespace string) map[string]types.AuthConfig {

	secretsAuthConfig := make(map[string]types.AuthConfig)

	for i := range secrets {
		res, err := k8sAPI.KubernetesClient.CoreV1().Secrets(namespace).Get(context.Background(), secrets[i], metav1.GetOptions{})
		if err != nil {
			logger.L().Error("unable to get secret", helpers.String("secret name", secrets[i]), helpers.Error(err))
			continue
		}
		sec, err := secrethandling.ParseSecret(res, secrets[i])
		if err != nil {
			logger.L().Error("failed to pars secret", helpers.String("secret name", secrets[i]), helpers.Error(err))
			continue
		}
		secretsAuthConfig[secrets[i]] = *sec
	}

	return secretsAuthConfig
}

// DEPRECATED
// GetImageRegistryCredentials returns various credentials for images in the pod
// imageTag empty means returns all of the credentials for all images in pod spec containers
// pod.ObjectMeta.Namespace must be well setted
func GetImageRegistryCredentials(imageTag string, pod *corev1.Pod) (map[string]types.AuthConfig, error) {
	k8sAPI := k8sinterface.NewKubernetesApi()
	listSecret, _ := listPodImagePullSecrets(&pod.Spec)
	listServiceSecret, _ := listServiceAccountImagePullSecrets(k8sAPI, pod.GetNamespace(), pod.Spec.ServiceAccountName)
	listSecret = append(listSecret, listServiceSecret...)
	secrets := getImagePullSecret(k8sAPI, listSecret, pod.ObjectMeta.Namespace)

	if len(secrets) == 0 {
		secrets = make(map[string]types.AuthConfig)
	}

	if imageTag != "" {
		cloudVendorSecrets, err := GetCloudVendorRegistryCredentials(imageTag)
		if err != nil {
			logger.L().Debug("failed to GetCloudVendorRegistryCredentials", helpers.String("imageTag", imageTag), helpers.Error(err))
		} else if len(cloudVendorSecrets) > 0 {
			for secName := range cloudVendorSecrets {
				secrets[secName] = cloudVendorSecrets[secName]
			}
		}
	} else {
		for contIdx := range pod.Spec.Containers {
			imageTag := pod.Spec.Containers[contIdx].Image

			cloudVendorSecrets, err := GetCloudVendorRegistryCredentials(imageTag)
			if err != nil {
				logger.L().Debug("failed to GetCloudVendorRegistryCredentials", helpers.String("imageTag", imageTag), helpers.Error(err))
			} else if len(cloudVendorSecrets) > 0 {
				for secName := range cloudVendorSecrets {
					secrets[secName] = cloudVendorSecrets[secName]
				}
			}
		}
	}

	return secrets, nil
}



// GetWorkloadsImages returns map[<image name>]<container name>
func GetWorkloadsImages(workload k8sinterface.IWorkload) map[string]string {
	images := map[string]string{}

	containers, err := workload.GetContainers()
	if err != nil {
		return images
	}
	for contIdx := range containers {
		images[containers[contIdx].Image] = containers[contIdx].Name
	}
	initContainers, err := workload.GetInitContainers()
	if err != nil {
		return images
	}
	for contIdx := range initContainers {
		images[initContainers[contIdx].Image] = initContainers[contIdx].Name
	}
	return images
}
