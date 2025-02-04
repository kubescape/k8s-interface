package cloudsupport

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/registry"
	"github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	"github.com/kubescape/k8s-interface/k8sinterface"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listPodImagePullSecrets(podSpec *corev1.PodSpec) ([]string, error) {
	if podSpec == nil {
		return []string{}, fmt.Errorf("in listPodImagePullSecrets podSpec is nil")
	}
	var secrets []string
	for _, i := range podSpec.ImagePullSecrets {
		secrets = append(secrets, i.Name)
	}
	return secrets, nil
}

func listServiceAccountImagePullSecrets(k8sAPI *k8sinterface.KubernetesApi, namespace, serviceAccountName string) ([]string, error) {
	var secrets []string
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

func getImagePullSecret(k8sAPI *k8sinterface.KubernetesApi, secrets []string, namespace string) map[string][]registry.AuthConfig {

	secretsAuthConfig := make(map[string][]registry.AuthConfig)

	for i := range secrets {
		res, err := k8sAPI.KubernetesClient.CoreV1().Secrets(namespace).Get(context.Background(), secrets[i], metav1.GetOptions{})
		if err != nil {
			logger.L().Error("unable to get secret", helpers.String("secret name", secrets[i]), helpers.Error(err))
			continue
		}
		sec, err := ParseSecret(res)
		if err != nil {
			logger.L().Error("failed to pars secret", helpers.String("secret name", secrets[i]), helpers.Error(err))
			continue
		}
		secretsAuthConfig[secrets[i]] = sec
	}

	return secretsAuthConfig
}

// GetImageRegistryCredentials returns various credentials for images in the pod
// imageTag empty means returns all of the credentials for all images in pod spec containers
// pod.ObjectMeta.Namespace must be well setted
// DEPRECATED
func GetImageRegistryCredentials(k8sAPI *k8sinterface.KubernetesApi, imageTag string, pod *corev1.Pod) (map[string][]registry.AuthConfig, error) {
	listSecret, _ := listPodImagePullSecrets(&pod.Spec)
	listServiceSecret, _ := listServiceAccountImagePullSecrets(k8sAPI, pod.GetNamespace(), pod.Spec.ServiceAccountName)
	listSecret = append(listSecret, listServiceSecret...)
	secrets := getImagePullSecret(k8sAPI, listSecret, pod.ObjectMeta.Namespace)

	if imageTag != "" {
		cloudVendorSecrets, err := GetCloudVendorRegistryCredentials(imageTag)
		if err != nil {
			logger.L().Debug("failed to GetCloudVendorRegistryCredentials", helpers.String("imageTag", imageTag), helpers.Error(err))
		} else if len(cloudVendorSecrets) > 0 {
			for secName := range cloudVendorSecrets {
				secrets[secName] = []registry.AuthConfig{cloudVendorSecrets[secName]}
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
					secrets[secName] = []registry.AuthConfig{cloudVendorSecrets[secName]}
				}
			}
		}
	}

	return secrets, nil
}

// GetWorkloadImageRegistryCredentials returns various credentials for images in the pod
// imageTag empty means returns all of the credentials for all images in pod spec containers
// pod.ObjectMeta.Namespace must be well setted
func GetWorkloadImageRegistryCredentials(k8sAPI *k8sinterface.KubernetesApi, imageTag string, workload k8sinterface.IWorkload) (map[string][]registry.AuthConfig, error) {
	podSpec, err := workload.GetPodSpec()
	if err != nil {
		return nil, err
	}
	listSecret, _ := listPodImagePullSecrets(podSpec)
	listServiceSecret, _ := listServiceAccountImagePullSecrets(k8sAPI, workload.GetNamespace(), podSpec.ServiceAccountName)
	listSecret = append(listSecret, listServiceSecret...)
	secrets := getImagePullSecret(k8sAPI, listSecret, workload.GetNamespace())

	if imageTag != "" {
		if cloudVendorSecrets, err := GetCloudVendorRegistryCredentials(imageTag); err != nil {
			return secrets, fmt.Errorf("failed to GetCloudVendorRegistryCredentials, image: %s, message: %v", imageTag, err)
		} else if len(cloudVendorSecrets) > 0 {
			for secName := range cloudVendorSecrets {
				secrets[secName] = []registry.AuthConfig{cloudVendorSecrets[secName]}
			}
		}
	} else {
		images := GetWorkloadsImages(workload)
		for imageTag := range images {
			if cloudVendorSecrets, err := GetCloudVendorRegistryCredentials(imageTag); err != nil {
				return secrets, fmt.Errorf("failed to GetCloudVendorRegistryCredentials, image: %s, message: %v", imageTag, err)
			} else if len(cloudVendorSecrets) > 0 {
				for secName := range cloudVendorSecrets {
					secrets[secName] = []registry.AuthConfig{cloudVendorSecrets[secName]}
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

type DockerConfigJsonstructure map[string]map[string]registry.AuthConfig

func updateSecret(authConfig *registry.AuthConfig, serverAddress string) {
	if authConfig.ServerAddress == "" {
		authConfig.ServerAddress = serverAddress
	}
	if authConfig.Username == "" || authConfig.Password == "" {
		auth := authConfig.Auth
		decodedAuth, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			return
		}

		splittedAuth := strings.Split(string(decodedAuth), ":")
		if len(splittedAuth) == 2 {
			authConfig.Username = splittedAuth[0]
			authConfig.Password = splittedAuth[1]
		}
	}
	if authConfig.Auth == "" {
		auth := fmt.Sprintf("%s:%s", authConfig.Username, authConfig.Password)
		authConfig.Auth = base64.StdEncoding.EncodeToString([]byte(auth))
	}
}

func parseEncodedSecret(sec map[string][]byte) (string, string) {
	buser := sec[corev1.BasicAuthUsernameKey]
	bpsw := sec[corev1.BasicAuthPasswordKey]
	duser, _ := base64.StdEncoding.DecodeString(string(buser))
	dpsw, _ := base64.StdEncoding.DecodeString(string(bpsw))
	return string(duser), string(dpsw)

}

func parseDecodedSecret(sec map[string]string) (string, string) {
	user := sec[corev1.BasicAuthUsernameKey]
	psw := sec[corev1.BasicAuthPasswordKey]
	return user, psw

}

func GetSecretContent(secret *corev1.Secret) (interface{}, error) {

	// Secret types- https://github.com/kubernetes/kubernetes/blob/7693a1d5fe2a35b6e2e205f03ae9b3eddcdabc6b/pkg/apis/core/types.go#L4394-L4478
	switch secret.Type {
	case corev1.SecretTypeDockerConfigJson:
		sec := make(DockerConfigJsonstructure)
		if err := json.Unmarshal(secret.Data[corev1.DockerConfigJsonKey], &sec); err != nil {
			return nil, err
		}
		return sec, nil
	default:
		user, psw := "", ""
		if len(secret.Data) != 0 {
			user, psw = parseEncodedSecret(secret.Data)
		} else if len(secret.StringData) != 0 {
			userD, pswD := parseDecodedSecret(secret.StringData)
			if userD != "" {
				user = userD
			}
			if pswD != "" {
				psw = pswD
			}
		} else {
			return nil, fmt.Errorf("data not found in secret")
		}
		if user == "" || psw == "" {
			return nil, fmt.Errorf("username  or password not found")
		}

		return &registry.AuthConfig{Username: user, Password: psw}, nil
	}
}

func ParseSecret(res *corev1.Secret) ([]registry.AuthConfig, error) {

	// Read secret
	secret, err := GetSecretContent(res)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("secret not found")
	}
	sec, err := ReadSecret(secret)
	if err != nil {
		return sec, err
	}
	return sec, nil

}

func ReadSecret(secret interface{}) ([]registry.AuthConfig, error) {
	// Store secret based on it's structure
	var authConfig []registry.AuthConfig
	if sec, ok := secret.(*registry.AuthConfig); ok {
		return []registry.AuthConfig{*sec}, nil
	}
	if sec, ok := secret.(map[string]string); ok {
		return []registry.AuthConfig{{Username: sec["username"]}}, nil
	}
	if sec, ok := secret.(DockerConfigJsonstructure); ok {
		if _, k := sec["auths"]; !k {
			return authConfig, fmt.Errorf("cant find auths")
		}
		for serverAddress, auth := range sec["auths"] {
			updateSecret(&auth, serverAddress)
			authConfig = append(authConfig, auth)
		}
		return authConfig, nil
	}

	return authConfig, fmt.Errorf("cant find secret")
}
