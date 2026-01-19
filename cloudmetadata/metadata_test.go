package cloudmetadata

import (
	"context"
	"testing"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/kubescape/k8s-interface/utils"
	v1 "k8s.io/api/core/v1"

	_ "embed"
)

//go:embed testdata/aks.json
var aks []byte

//go:embed testdata/digitalocean.json
var digitalocean []byte

//go:embed testdata/eks.json
var eks []byte

//go:embed testdata/gke.json
var gke []byte

//go:embed testdata/linode.json
var linode []byte

//go:embed testdata/vsphere.json
var vsphere []byte

func TestGetCloudMetadata(t *testing.T) {
	tests := []struct {
		name    string
		node    []byte
		want    *armotypes.CloudMetadata
		wantErr bool
	}{
		{
			name: "AWS provider",
			node: eks,
			want: &armotypes.CloudMetadata{
				AccountID:    "",
				Provider:     armotypes.ProviderAws,
				InstanceType: "c7g.large",
				Region:       "eu-west-1",
				Zone:         "eu-west-1b",
				InstanceID:   "i-000000fff00eeaa00",
				PrivateIP:    "1.1.1.1",
				PublicIP:     "",
				Hostname:     "ip-1-1-1-1.eu-west-1.compute.internal",
			},
			wantErr: false,
		},
		{
			name: "DigitalOcean provider",
			node: digitalocean,
			want: &armotypes.CloudMetadata{
				AccountID:    "",
				Provider:     armotypes.ProviderDigitalOcean,
				InstanceType: "s-8vcpu-16gb",
				Region:       "fra1",
				Zone:         "",
				InstanceID:   "123456789",
				PrivateIP:    "1.1.1.2",
				PublicIP:     "1.1.1.6",
				Hostname:     "default-pool-test123",
			},
			wantErr: false,
		},
		{
			name: "Azure provider",
			node: aks,
			want: &armotypes.CloudMetadata{
				AccountID:    "00000000-ffff-ffff-ffff-fffffffffff4",
				Provider:     armotypes.ProviderAzure,
				InstanceType: "Standard_D8ads_v5",
				Region:       "westeurope",
				Zone:         "westeurope-1",
				InstanceID:   "aks-minio01-00000000-vmss",
				PrivateIP:    "1.1.1.1",
				PublicIP:     "",
				Hostname:     "aks-minio01-00000000-vmss00000k",
			},
			wantErr: false,
		},
		{
			name: "GCP provider",
			node: gke,
			want: &armotypes.CloudMetadata{
				AccountID:    "kubescape-123456",
				Provider:     armotypes.ProviderGcp,
				InstanceType: "c2-standard-16",
				Region:       "us-west1",
				Zone:         "us-west1-a",
				InstanceID:   "gke-cluster-pool-1-123456",
				PrivateIP:    "1.1.1.2",
				PublicIP:     "",
				Hostname:     "gke-cluster-pool-1-123456",
			},
			wantErr: false,
		},
		{
			name: "Linode provider",
			node: linode,
			want: &armotypes.CloudMetadata{
				Provider:     armotypes.ProviderLinode,
				InstanceType: "g6-standard-1",
				Region:       "eu-central",
				InstanceID:   "71504446",
				PrivateIP:    "192.168.137.215",
				Hostname:     "lke99535-149699-641dcb2a5313",
			},
			wantErr: false,
		},
		{
			name: "VMware provider",
			node: vsphere,
			want: &armotypes.CloudMetadata{
				AccountID:    "",
				Provider:     armotypes.ProviderVMware,
				InstanceType: "",
				Region:       "",
				Zone:         "",
				InstanceID:   "00000000-ffff-ffff-ffff-fffffffffff9",
				PrivateIP:    "1.1.1.1",
				PublicIP:     "",
				Hostname:     "aaaa00",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := utils.ConvertUnstructuredToRuntimeObject(utils.MustConvertRawToUnstructured(tt.node))
			if err != nil {
				t.Fatalf("ConvertUnstructuredToRuntimeObject() error = %v", err)
			}
			nodeObj := n.(*v1.Node)
			got, err := GetCloudMetadata(context.Background(), nodeObj, nodeObj.GetName())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCloudMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareCloudMetadata(got, tt.want) {
				t.Errorf("GetCloudMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareCloudMetadata(a, b *armotypes.CloudMetadata) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Provider == b.Provider &&
		a.InstanceType == b.InstanceType &&
		a.Region == b.Region &&
		a.Zone == b.Zone &&
		a.InstanceID == b.InstanceID &&
		a.PrivateIP == b.PrivateIP &&
		a.PublicIP == b.PublicIP &&
		a.Hostname == b.Hostname &&
		a.AccountID == b.AccountID
}
