## List of environment variables related to cloud

### General

| Variable | Description |
|---|---|
| `KS_KUBE_CLUSTER` | Cluster name used for cloud provider lookups instead of the current kubeconfig context name (`cloudsupport.GetKubeContextName`). |
| `KS_OFFLINE` | When set to `true`, the `Get*FromCloudProvider` functions in `cloudsupport` skip all cloud provider API calls and return `ErrCloudDescribeUnavailable`, so callers can continue without cloud data (for example in air-gapped clusters). The Kubescape Helm chart sets it when `capabilities.kubescapeOffline` is `enable`. |

### AKS

| Variable | Description |
|---|---|
| `AZURE_SUBSCRIPTION_ID` | Azure subscription ID of the cluster. Required for AKS lookups. |
| `AZURE_RESOURCE_GROUP` | Azure resource group of the cluster. Required for AKS lookups. |
| `AZURE_PULL_CLIENT_ID` | Client ID (UUID) of the managed identity used to pull images from ACR. When unset, the managed identity endpoint uses its default identity. |

### EKS

| Variable | Description |
|---|---|
| `KS_CLOUD_REGION` | AWS region of the cluster. |
| `AWS_REGION` | Used as the region when `KS_CLOUD_REGION` is not set. If neither is set, the region is taken from the default AWS SDK configuration, then parsed from the cluster name. |
| `KS_EKS_INVENTORY_TIMEOUT` | Optional overall deadline for the IAM inventory (`GetListEntitiesForPolicies` and `GetPolicyVersion`: all `ListPolicies` pages plus one call per policy), as a Go duration such as `2m`. When unset there is no overall deadline; each AWS request is still bounded by its own timeout. A value that does not parse or is not positive (for example `30` with no unit) is logged as a warning and ignored. |
| `AWS_EC2_METADATA_SERVICE_ENDPOINT` | Overrides the EC2 instance metadata service (IMDS) endpoint used by the `cloudmetadata` AWS metadata client, matching the AWS SDK. Defaults to `http://169.254.169.254`. |

### GKE

| Variable | Description |
|---|---|
| `KS_GKE_PROJECT` | GCP project of the cluster. When unset, it is parsed from a `gke_<project>_<region>_<cluster>` context name. |
| `KS_CLOUD_REGION` | GCP region or zone of the cluster. When unset or empty, it is parsed from a `gke_<project>_<region>_<cluster>` context name. |
| `GOOGLE_CLOUD_QUOTA_PROJECT` | Read by the Google Cloud client libraries. When set, it takes precedence over the `quota_project_id` of the Application Default Credentials file. |
