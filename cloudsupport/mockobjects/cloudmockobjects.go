package mockobjects

var EksDescriptor = `{
    "Cluster": {
        "Arn": "arn:aws:eks:eu-west-23:XXXXXXXXXXXX:cluster/my-cluster",
        "CertificateAuthority": {
            "Data": ""
        },
        "ClientRequestToken": null,
        "ConnectorConfig": null,
        "EncryptionConfig": null,
        "Endpoint": "https://XXXXXXXXXXXX.r45.eu-west-23.eks.amazonaws.com",
        "Identity": {
            "Oidc": {
                "Issuer": "https://oidc.eks.eu-west-23.amazonaws.com/id/XXXXXXXXXXXX"
            }
        },
        "KubernetesNetworkConfig": {
            "ServiceIpv4Cidr": "0.0.0.0/16"
        },
        "Logging": {
            "ClusterLogging": [
                {
                    "Enabled": false,
                    "Types": [
                        "api",
                        "audit",
                        "authenticator",
                        "controllerManager",
                        "scheduler"
                    ]
                }
            ]
        },
        "Name": "ca-terraform-eks-dev-stage",
        "PlatformVersion": "eks.0",
        "ResourcesVpcConfig": {
            "ClusterSecurityGroupId": "sg-XXXXXXXXXXXX",
            "EndpointPrivateAccess": true,
            "EndpointPublicAccess": true,
            "PublicAccessCidrs": [
                "0.0.0.0/0"
            ],
            "SecurityGroupIds": [
                "sg-XXXXXXXXXXXX"
            ],
            "SubnetIds": [
                "subnet-XXXXXXXXXXX0",
                "subnet-XXXXXXXXXXX1",
                "subnet-XXXXXXXXXXX2"
            ],
            "VpcId": "vpc-XXXXXXXXXXXX"
        },
        "RoleArn": "arn:aws:iam::XXXXXXXXXXXX:role/terraform-XXXXXXXXXXXX",
        "Status": "ACTIVE",
        "Tags": {
            "Customer": "Armo",
            "Name": "terraform-eks-XXXXXXXXXXXX",
            "Owner": "my-cluster",
            "Project": "Infra"
        },
        "Version": "0.0"
    }
}
`

var AksDescriptor = `     {
    "identity": {
        "type": "SystemAssigned"
    },
    "location": "westeurope",
    "properties": {
        "addonProfiles": {
            "azureKeyvaultSecretsProvider": {
                "enabled": false
            },
            "azurepolicy": {
                "enabled": false
            },
            "httpApplicationRouting": {
                "enabled": false
            }
        },
        "agentPoolProfiles": [{
            "availabilityZones": ["1", "2", "3"],
            "count": 1,
            "enableAutoScaling": true,
            "enableFIPS": false,
            "enableNodePublicIP": false,
            "kubeletDiskType": "OS",
            "maxCount": 5,
            "maxPods": 110,
            "minCount": 1,
            "mode": "System",
            "name": "agentpool",
            "orchestratorVersion": "1.21.9",
            "osDiskSizeGB": 128,
            "osDiskType": "Managed",
            "osSKU": "Ubuntu",
            "osType": "Linux",
            "powerState": {
                "code": "Running"
            },
            "tags": {
                "Owner": "Daniel"
            },
            "type": "VirtualMachineScaleSets",
            "vmSize": "Standard_B2s"
        }],
        "apiServerAccessProfile": {
            "enablePrivateCluster": false
        },
        "autoScalerProfile": {
            "balance-similar-node-groups": "false",
            "expander": "random",
            "max-empty-bulk-delete": "10",
            "max-graceful-termination-sec": "600",
            "max-node-provision-time": "15m",
            "max-total-unready-percentage": "45",
            "new-pod-scale-up-delay": "0s",
            "ok-total-unready-count": "3",
            "scan-interval": "10s",
            "scale-down-delay-after-add": "10m",
            "scale-down-delay-after-delete": "10s",
            "scale-down-delay-after-failure": "3m",
            "scale-down-unneeded-time": "10m",
            "scale-down-unready-time": "20m",
            "scale-down-utilization-threshold": "0.5",
            "skip-nodes-with-local-storage": "false",
            "skip-nodes-with-system-pods": "true"
        },
        "dnsPrefix": "armo-testing",
        "enableRBAC": true,
        "identityProfile": {
            "kubeletidentity": {
                "resourceId": "/subscriptions/XXXXX/resourcegroups/MC_armo-dev_armo-testing_westeurope/providers/Microsoft.ManagedIdentity/userAssignedIdentities/armo-testing-agentpool",
                "clientId": "XXXXX",
                "objectId": "XXXXX"
            }
        },
        "kubernetesVersion": "1.21.9",
        "networkProfile": {
            "networkPlugin": "kubenet",
            "podCidr": "10.244.0.0/16",
            "serviceCidr": "10.0.0.0/16",
            "dnsServiceIP": "10.0.0.10",
            "dockerBridgeCidr": "172.17.0.1/16",
            "outboundType": "loadBalancer",
            "loadBalancerSku": "Standard",
            "loadBalancerProfile": {
                "managedOutboundIPs": {
                    "count": 1
                },
                "effectiveOutboundIPs": [{
                    "id": "XXXX"
                }]
            }
        },
        "nodeResourceGroup": "MC_armo-dev_armo-testing_westeurope",
        "servicePrincipalProfile": {
            "clientId": "msi"
        }
    },
    "sku": {
        "name": "Basic",
        "tier": "Free"
    },
    "tags": {
        "Owner": "XXXX"
    }
}`

var GkeDescriptor = `    {
	"name": "kubescape-demo-01",
	"node_config": {
		"machine_type": "e2-medium",
		"disk_size_gb": 100,
		"oauth_scopes": [
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/https://console.cloud.google.com/kubernetes/clusters/details/us-central1-c/kubescape-demo-01/details?authuser=0&project=kubescape-demo-01/monitoring",
			"https://www.googleapis.com/auth/servicecontrol",
			"https://www.googleapis.com/auth/service.management.readonly",
			"https://www.googleapis.com/auth/trace.append"
		],
		"service_account": "default",
		"metadata": {
			"disable-legacy-endpoints": "true"
		},
		"image_type": "COS_CONTAINERD",
		"disk_type": "pd-standard",
		"shielded_instance_config": {
			"enable_integrity_monitoring": true
		}
	},
	"master_auth": {
		"cluster_ca_certificate": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	},
	"logging_service": "logging.googleapis.com/kubernetes",
	"monitoring_service": "monitoring.googleapis.com/kubernetes",
	"network": "default",
	"cluster_ipv4_cidr": "0.0.0.0/14",
	"addons_config": {
		"http_load_balancing": {},
		"horizontal_pod_autoscaling": {},
		"kubernetes_dashboard": {
			"disabled": true
		},
		"network_policy_config": {
			"disabled": true
		},
		"dns_cache_config": {}
	},
	"subnetwork": "default",
	"node_pools": [
		{
			"name": "default-pool",
			"config": {
				"machine_type": "e2-medium",
				"disk_size_gb": 100,
				"oauth_scopes": [
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring",
					"https://www.googleapis.com/auth/servicecontrol",
					"https://www.googleapis.com/auth/service.management.readonly",
					"https://www.googleapis.com/auth/trace.append"
				],
				"service_account": "default",
				"metadata": {
					"disable-legacy-endpoints": "true"
				},
				"image_type": "COS_CONTAINERD",
				"disk_type": "pd-standard",
				"shielded_instance_config": {
					"enable_integrity_monitoring": true
				}
			},
			"initial_node_count": 3,
			"locations": [
				"us-central1-r"
			],
			"self_link": "https://container.googleapis.com/v1/projects/kubescape-demo-01/zones/us-central1-r/clusters/kubescape-demo-01/nodePools/default-pool",
			"version": "0.0.0-gke.0",
			"instance_group_urls": [
				"https://www.googleapis.com/compute/v1/projects/kubescape-demo-01/zones/us-central1-r/instanceGroupManagers/kubescape-demo-01-grp"
			],
			"status": 2,
			"autoscaling": {},
			"management": {
				"auto_upgrade": true,
				"auto_repair": true
			},
			"max_pods_constraint": {
				"max_pods_per_node": 110
			},
			"pod_ipv4_cidr_size": 24,
			"upgrade_settings": {
				"max_surge": 1
			}
		}
	],
	"locations": [
		"us-central1-r"
	],
	"label_fingerprint": "dfsdfsdf",
	"legacy_abac": {},
	"ip_allocation_policy": {
		"use_ip_aliases": true,
		"cluster_ipv4_cidr": "0.0.0.0/14",
		"services_ipv4_cidr": "0.0.0.0/20",
		"cluster_secondary_range_name": "kubescape-demo-01-686ce31a",
		"services_secondary_range_name": "kubescape-demo-01-686ce31a",
		"cluster_ipv4_cidr_block": "0.0.0.0/14",
		"services_ipv4_cidr_block": "0.0.0.0/20"
	},
	"master_authorized_networks_config": {},
	"maintenance_policy": {
		"resource_version": "1651165"
	},
	"autoscaling": {},
	"network_config": {
		"network": "projects/kubescape-demo-01/global/networks/default",
		"subnetwork": "projects/kubescape-demo-01/regions/us-central1-r/subnetworks/default",
		"default_snat_status": {}
	},
	"default_max_pods_constraint": {
		"max_pods_per_node": 110
	},
	"authenticator_groups_config": {},
	"database_encryption": {
		"state": 2
	},
	"shielded_nodes": {
		"enabled": true
	},
	"release_channel": {
		"channel": 2
	},
	"self_link": "https://container.googleapis.com/v1/projects/kubescape-demo-01/zones/us-central1-r/clusters/kubescape-demo-01",
	"zone": "us-central1-r",
	"endpoint": "0.0.0.0",
	"initial_cluster_version": "0-gke.0",
	"current_master_version": "0-gke.0",
	"current_node_version": "0-gke.0",
	"status": 2,
	"services_ipv4_cidr": "0.0.0.0/20",
	"instance_group_urls": [
		"https://www.googleapis.com/compute/v1/projects/kubescape-demo-01/zones/us-central1-r/instanceGroupManagers/kubescape-demo-01-grp"
	],
	"current_node_count": 3,
	"location": "us-central1-r"
}`

var EksDescribeRepositories = `{
        "NextToken": null,
        "Repositories": [
            {
                "CreatedAt": "2023-01-11T13:07:39Z",
                "EncryptionConfiguration": {
                    "EncryptionType": "XXXXXX",
                    "KmsKey": null
                },
                "ImageScanningConfiguration": {
                    "ScanOnPush": true
                },
                "ImageTagMutability": "MUTABLE",
                "RegistryId": "XXXXXXXXXXXX",
                "RepositoryArn": "arn:aws:ecr:eu-west-2:XXXXXXXXXXXX:repository/my-repo",
                "RepositoryName": "my-repo",
                "RepositoryUri": "XXXXXXXXXXXX.dkr.ecr.eu-west-2.amazonaws.com/my-repo"
            },
            {
                "CreatedAt": "2023-01-11T13:07:39Z",
                "EncryptionConfiguration": {
                    "EncryptionType": "XXXXXX",
                    "KmsKey": null
                },
                "ImageTagMutability": "MUTABLE",
                "RegistryId": "XXXXXXXXXXXX",
                "RepositoryArn": "arn:aws:ecr:eu-west-2:XXXXXXXXXXXX:repository/my-repo",
                "RepositoryName": "my-repo-no-scanning",
                "RepositoryUri": "XXXXXXXXXXXX.dkr.ecr.eu-west-2.amazonaws.com/my-repo"
            }
        ],
        "ResultMetadata": {}
}
    `

var EksListEntitiesForPolicies = `
{
    "rolesPolicies": {
        "arn:aws:iam::XXXXXXXXXXXX:policy/FullAccess": {
            "IsTruncated": false,
            "Marker": null,
            "PolicyGroups": [
                {
                    "GroupId": "XXXXXXXXXXXX",
                    "GroupName": "Developers"
                }
            ],
            "PolicyRoles": [],
            "PolicyUsers": [
                {
                    "UserId": "XXXXXXXXXXXX",
                    "UserName": "my-user"
                }
            ],
            "ResultMetadata": {}
        },
        "arn:aws:iam::XXXXXXXXXXXX:policy/DevStage": {
            "IsTruncated": false,
            "Marker": null,
            "PolicyGroups": [],
            "PolicyRoles": [
                {
                    "RoleId": "XXXXXXXXXXXX",
                    "RoleName": "eks-dev"
                },
                {
                    "RoleId": "XXXXXXXXXXXX",
                    "RoleName": "eks-dev-test"
                }
            ],
            "PolicyUsers": [],
            "ResultMetadata": {}
        }
    }
}
        `
var EksGetPolicyVersion = `
{
    "apiVersion": "eks.amazonaws.com/v1",
    "kind": "PolicyVersion",
    "metadata": {
        "name": "ca-terraform-eks-dev-stage",
        "provider": "eks"
    },
    "data": {
        "policiesDocuments": {
            "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly": {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Action": [
                            "ecr:GetAuthorizationToken",
                            "ecr:BatchCheckLayerAvailability",
                            "ecr:GetDownloadUrlForLayer",
                            "ecr:GetRepositoryPolicy",
                            "ecr:DescribeRepositories",
                            "ecr:ListImages",
                            "ecr:DescribeImages",
                            "ecr:BatchGetImage",
                            "ecr:GetLifecyclePolicy",
                            "ecr:GetLifecyclePolicyPreview",
                            "ecr:ListTagsForResource",
                            "ecr:DescribeImageScanFindings"
                        ],
                        "Resource": "*"
                    }
                ]
            },
            "arn:aws:iam::aws:policy/AWSMarketplaceFullAccess": {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Action": "ec2:*",
                        "Effect": "Allow",
                        "Resource": "*"
                    },
                    {
                        "Effect": "Allow",
                        "Action": "elasticloadbalancing:*",
                        "Resource": "*"
                    }
                ]
            }
        }
    }
}
`

var AKSListRoleAssignments = `
    {
        "roleAssignments": [
            {
                "id": "/subscriptions/********/providers/Microsoft.Authorization/roleAssignments/b925e22e-a17f-4228-b268-bed5500e463a",
                "name": "b925e22e-a17f-4228-b268-bed5500e463a",
                "properties": {
                    "createdBy": "*******",
                    "createdOn": "2022-12-25T09:58:55.3897463Z",
                    "principalId": "123456789",
                    "principalType": "User",
                    "roleDefinitionId": "/subscriptions/********/providers/Microsoft.Authorization/roleDefinitions/2222222",
                    "scope": "/subscriptions/********",
                    "updatedBy": "*******",
                    "updatedOn": "2022-12-25T09:58:55.3897463Z"
                },
                "type": "Microsoft.Authorization/roleAssignments"
            },
            {
                "id": "/subscriptions/*********/providers/Microsoft.Authorization/roleAssignments/98765",
                "name": "98765",
                "properties": {
                    "createdBy": "******",
                    "createdOn": "2023-02-12T11:50:35.2541275Z",
                    "principalId": "8fd8eace-5423-4110-b6c8-38f94223bca7",
                    "principalType": "Group",
                    "roleDefinitionId": "/subscriptions/*********/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
                    "scope": "/subscriptions/*********",
                    "updatedBy": "******",
                    "updatedOn": "2023-02-12T11:50:35.2541275Z"
                },
                "type": "Microsoft.Authorization/roleAssignments"
            },
            {
                "id": "/subscriptions/*********/resourceGroups/MC_matan-nypnk-rg_matan-nypnk-aks_northeurope/providers/Microsoft.Authorization/roleAssignments/987654321",
                "name": "987654321",
                "properties": {
                    "createdBy": "123456",
                    "createdOn": "2023-03-14T11:49:41.0899644Z",
                    "principalId": "1234561234",
                    "principalType": "ServicePrincipal",
                    "roleDefinitionId": "/subscriptions/*********/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
                    "scope": "/subscriptions/*********/resourceGroups/MC_matan-nypnk-rg_matan-nypnk-aks_northeurope",
                    "updatedBy": "123456",
                    "updatedOn": "2023-03-14T11:49:41.0899644Z"
                },
                "type": "Microsoft.Authorization/roleAssignments"
            }
        ]
}`

var AKSListRoleDefinitions = `
{
    "roleDefinitions": [
        {
            "properties": {
                "assignableScopes": [
                    "/"
                ],
                "description": "acr push",
                "permissions": [
                    {
                        "actions": [
                            "Microsoft.ContainerRegistry/registries/pull/read",
                            "Microsoft.ContainerRegistry/registries/push/write"
                        ],
                        "notActions": []
                    }
                ],
                "roleName": "AcrPush",
                "type": "BuiltInRole"
            },
            "id": "/subscriptions/05fed2cc-041d-4c53-b4b8-a1aa128d189f/providers/Microsoft.Authorization/roleDefinitions/8311e382-0749-4cb8-b61a-304f252e45ec",
            "name": "8311e382-0749-4cb8-b61a-304f252e45ec",
            "type": "Microsoft.Authorization/roleDefinitions"
        },
        {
            "properties": {
                "assignableScopes": [
                    "/"
                ],
                "description": "Used by the Avere vFXT cluster to manage the cluster",
                "permissions": [
                    {
                        "actions": [
                            "Microsoft.Compute/virtualMachines/read",
                            "Microsoft.Network/networkInterfaces/read",
                            "Microsoft.Network/virtualNetworks/read",
                            "Microsoft.Network/virtualNetworks/subnets/read",
                            "Microsoft.Resources/subscriptions/resourceGroups/read",
                            "Microsoft.Storage/storageAccounts/blobServices/containers/read"
                        ],
                        "notActions": []
                    }
                ],
                "roleName": "Avere Operator",
                "type": "BuiltInRole"
            },
            "id": "/subscriptions/05fed2cc-041d-4c53-b4b8-a1aa128d189f/providers/Microsoft.Authorization/roleDefinitions/c025889f-8102-4ebf-b32c-fc0c6f0c6bd9",
            "name": "c025889f-8102-4ebf-b32c-fc0c6f0c6bd9",
            "type": "Microsoft.Authorization/roleDefinitions"
        },
        {
            "properties": {
                "assignableScopes": [
                    "/"
                ],
                "description": "Lets you manage Azure Stack registrations.",
                "permissions": [
                    {
                        "actions": [
                            "Microsoft.AzureStack/edgeSubscriptions/read",
                            "Microsoft.AzureStack/registrations/products/*/action",
                            "Microsoft.AzureStack/registrations/products/read",
                            "Microsoft.AzureStack/registrations/read"
                        ],
                        "notActions": []
                    }
                ],
                "roleName": "Azure Stack Registration Owner",
                "type": "BuiltInRole"
            },
            "id": "/subscriptions/05fed2cc-041d-4c53-b4b8-a1aa128d189f/providers/Microsoft.Authorization/roleDefinitions/6f12a6df-dd06-4f3e-bcb1-ce8be600526a",
            "name": "6f12a6df-dd06-4f3e-bcb1-ce8be600526a",
            "type": "Microsoft.Authorization/roleDefinitions"
        }
    ]
}
`
