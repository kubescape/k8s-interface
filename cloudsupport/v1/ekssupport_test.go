package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kubescape/k8s-interface/k8sinterface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetContextName(t *testing.T) {
	defer tearDown()

	mockname1 := "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster"
	eksSupport := NewEKSSupport()
	name := eksSupport.GetContextName(mockname1)
	assert.Equal(t, "test-cluster", name)
	region, err := eksSupport.GetRegion(mockname1)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	mockname2 := "arn:aws:eks:eu-north-1:123456789:cluster/test-cluster"
	splittedCluster := strings.Split(mockname2, "/")
	name = splittedCluster[len(splittedCluster)-1]
	assert.Equal(t, "test-cluster", name)
	region, err = eksSupport.GetRegion(mockname2)
	assert.NoError(t, err)
	assert.Equal(t, "eu-north-1", region)

	tests := []struct {
		name                string
		config              *clientcmdapi.Config
		connected           bool
		cluster             string
		expectedContextName string
	}{
		{
			name: "No heuristic matches, falls back to explicit cluster name",
			config: &clientcmdapi.Config{
				CurrentContext: "d34db33f",
				Clusters: map[string]*clientcmdapi.Cluster{
					"d34db33f": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           true,
			cluster:             "my-cluster",
			expectedContextName: "my-cluster",
		},
		{
			name: "Context name is cluster name, connected",
			config: &clientcmdapi.Config{
				CurrentContext: "my-cluster",
				Clusters: map[string]*clientcmdapi.Cluster{
					"my-cluster": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           true,
			cluster:             "my-cluster",
			expectedContextName: "my-cluster",
		},
		{
			name: "Context name is cluster name, not connected",
			config: &clientcmdapi.Config{
				CurrentContext: "my-cluster",
				Clusters: map[string]*clientcmdapi.Cluster{
					"my-cluster": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           false,
			cluster:             "my-cluster",
			expectedContextName: "my-cluster",
		},
		{
			name: "No cluster and no matching context, returns empty string",
			config: &clientcmdapi.Config{
				CurrentContext: "d34db33f",
				Clusters: map[string]*clientcmdapi.Cluster{
					"d34db33f": {
						Server: "https://my-server.local",
					},
				},
			},
			connected:           true,
			cluster:             "",
			expectedContextName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sinterface.SetConnectedToCluster(tt.connected)
			k8sinterface.SetClientConfigAPI(tt.config)
			actualContextName := eksSupport.GetContextName(tt.cluster)
			assert.Equal(t, tt.expectedContextName, actualContextName)
		})
	}
}

func TestGetRegion(t *testing.T) {
	tests := []struct {
		name           string
		cluster        string
		envRegion      string
		expectedRegion string
		expectErr      bool
	}{
		{
			name:           "Region is extracted from cluster name 1",
			cluster:        "arn:aws:eks:eu-north-1:123456789:cluster-test-cluster",
			expectedRegion: "eu-north-1",
			expectErr:      false,
		},
		{
			name:           "Region is extracted from cluster name 2",
			cluster:        "arn-aws-eks-eu-west-2-XXXXXXXXXXXX-cluster-Yiscah-test-g2am5",
			expectedRegion: "eu-west-2",
			expectErr:      false,
		},
		{
			name:           "Region is present in environment variable",
			envRegion:      "us-west-2",
			expectedRegion: "us-west-2",
			expectErr:      false,
		},
		{
			name:           "Region is extracted from cluster name with ':' separator",
			cluster:        "cluster:us-west-2:eks",
			expectedRegion: "us-west-2",
			expectErr:      true,
		},
		{
			name:      "Failed to get region",
			cluster:   "cluster",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envRegion != "" {
				os.Setenv(KS_CLOUD_REGION_ENV_VAR, tt.envRegion)
				defer os.Unsetenv(KS_CLOUD_REGION_ENV_VAR)
			}

			eksSupport := &EKSSupport{}
			region, err := eksSupport.GetRegion(tt.cluster)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRegion, region)
			}
		})
	}

	t.Run("AWS_REGION environment variable is checked", func(t *testing.T) {
		awsRegion := "ap-southeast-2"
		os.Setenv("AWS_REGION", awsRegion)
		defer os.Unsetenv("AWS_REGION")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, awsRegion, region)
	})

	t.Run("KS_CLOUD_REGION takes precedence over AWS_REGION", func(t *testing.T) {
		ksRegion := "us-west-2"
		awsRegion := "eu-west-1"
		os.Setenv(KS_CLOUD_REGION_ENV_VAR, ksRegion)
		os.Setenv("AWS_REGION", awsRegion)
		defer os.Unsetenv(KS_CLOUD_REGION_ENV_VAR)
		defer os.Unsetenv("AWS_REGION")

		eksSupport := &EKSSupport{}
		region, err := eksSupport.GetRegion("my-cluster")

		assert.NoError(t, err)
		assert.Equal(t, ksRegion, region)
	})
}

func TestEKSTimeoutConstants(t *testing.T) {
	assert.Equal(t, 5*time.Second, eksCallTimeout)
	assert.Equal(t, 30*time.Second, eksRBACEnumerationTimeout)
	assert.Less(t, eksCallTimeout, eksRBACEnumerationTimeout)
}

func TestGetClusterDescribeTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(eksCallTimeout + 500*time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"cluster":{"name":"test-cluster","status":"ACTIVE"}}`))
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	eksSupport := NewEKSSupport()
	start := time.Now()
	_, err := eksSupport.GetClusterDescribe("test-cluster", "us-east-1")
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.GreaterOrEqual(t, elapsed, eksCallTimeout)
	assert.Less(t, elapsed, eksCallTimeout+2*time.Second)
}

func TestGetListEntitiesForPolicies_Behavioral(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		action := r.Form.Get("Action")
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		switch action {
		case "ListPolicies":
			w.Write([]byte(`
<ListPoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListPoliciesResult>
    <Policies>
      <member>
        <Arn>arn:aws:iam::123456789012:policy/PolicyA</Arn>
        <PolicyName>PolicyA</PolicyName>
        <DefaultVersionId>v1</DefaultVersionId>
      </member>
      <member>
        <Arn>arn:aws:iam::123456789012:policy/PolicyB</Arn>
        <PolicyName>PolicyB</PolicyName>
        <DefaultVersionId>v1</DefaultVersionId>
      </member>
    </Policies>
    <IsTruncated>false</IsTruncated>
  </ListPoliciesResult>
  <ResponseMetadata><RequestId>req-1</RequestId></ResponseMetadata>
</ListPoliciesResponse>`))
		case "ListEntitiesForPolicy":
			arn := r.Form.Get("PolicyArn")
			w.Write([]byte(fmt.Sprintf(`
<ListEntitiesForPolicyResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListEntitiesForPolicyResult>
    <PolicyGroups/>
    <PolicyUsers/>
    <PolicyRoles>
      <member>
        <RoleId>role-1</RoleId>
        <RoleName>RoleFor-%s</RoleName>
        <Arn>arn:aws:iam::123456789012:role/RoleFor-%s</Arn>
      </member>
    </PolicyRoles>
    <IsTruncated>false</IsTruncated>
  </ListEntitiesForPolicyResult>
  <ResponseMetadata><RequestId>req-2</RequestId></ResponseMetadata>
</ListEntitiesForPolicyResponse>`, arn, arn)))
		default:
			w.Write([]byte(`<ErrorResponse><Error><Code>InvalidAction</Code></Error></ErrorResponse>`))
		}
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	eksSupport := NewEKSSupport()
	out, err := eksSupport.GetListEntitiesForPolicies("us-east-1")
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Len(t, out.EntitiesForPolicies, 2)
	assert.Contains(t, out.EntitiesForPolicies, "arn:aws:iam::123456789012:policy/PolicyA")
	assert.Contains(t, out.EntitiesForPolicies, "arn:aws:iam::123456789012:policy/PolicyB")
}

func TestGetPolicyVersion_Behavioral(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		action := r.Form.Get("Action")
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		switch action {
		case "ListPolicies":
			w.Write([]byte(`
<ListPoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListPoliciesResult>
    <Policies>
      <member>
        <Arn>arn:aws:iam::123456789012:policy/PolicyV</Arn>
        <PolicyName>PolicyV</PolicyName>
        <DefaultVersionId>v1</DefaultVersionId>
      </member>
    </Policies>
    <IsTruncated>false</IsTruncated>
  </ListPoliciesResult>
  <ResponseMetadata><RequestId>req-1</RequestId></ResponseMetadata>
</ListPoliciesResponse>`))
		case "GetPolicyVersion":
			w.Write([]byte(`
<GetPolicyVersionResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <GetPolicyVersionResult>
    <PolicyVersion>
      <Document>%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Action%22%3A%5B%22s3%3AListBucket%22%5D%2C%22Resource%22%3A%22%2A%22%7D%5D%7D</Document>
      <VersionId>v1</VersionId>
      <IsDefaultVersion>true</IsDefaultVersion>
    </PolicyVersion>
  </GetPolicyVersionResult>
  <ResponseMetadata><RequestId>req-3</RequestId></ResponseMetadata>
</GetPolicyVersionResponse>`))
		default:
			w.Write([]byte(`<ErrorResponse><Error><Code>InvalidAction</Code></Error></ErrorResponse>`))
		}
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	eksSupport := NewEKSSupport()
	out, err := eksSupport.GetPolicyVersion("us-east-1")
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Len(t, out.PolicyVersion, 1)
	pv, ok := out.PolicyVersion["arn:aws:iam::123456789012:policy/PolicyV"]
	require.True(t, ok)
	assert.Equal(t, "2012-10-17", pv.Version)
	require.Len(t, pv.Statement, 1)
	assert.Equal(t, "Allow", pv.Statement[0].Effect)
	assert.Equal(t, []string{"s3:ListBucket"}, pv.Statement[0].Action)
}

func TestGetListEntitiesForPolicies_ConfigurableTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		action := r.Form.Get("Action")
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		switch action {
		case "ListPolicies":
			w.Write([]byte(`
<ListPoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListPoliciesResult>
    <Policies>
      <member><Arn>arn:aws:iam::123456789012:policy/P1</Arn><PolicyName>P1</PolicyName><DefaultVersionId>v1</DefaultVersionId></member>
      <member><Arn>arn:aws:iam::123456789012:policy/P2</Arn><PolicyName>P2</PolicyName><DefaultVersionId>v1</DefaultVersionId></member>
    </Policies>
    <IsTruncated>false</IsTruncated>
  </ListPoliciesResult>
  <ResponseMetadata><RequestId>req-1</RequestId></ResponseMetadata>
</ListPoliciesResponse>`))
		case "ListEntitiesForPolicy":
			// Each call takes 150ms
			time.Sleep(150 * time.Millisecond)
			w.Write([]byte(`
<ListEntitiesForPolicyResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListEntitiesForPolicyResult><PolicyGroups/><PolicyUsers/><PolicyRoles/><IsTruncated>false</IsTruncated></ListEntitiesForPolicyResult>
  <ResponseMetadata><RequestId>req-2</RequestId></ResponseMetadata>
</ListEntitiesForPolicyResponse>`))
		}
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	eksSupport := NewEKSSupport()

	// When KS_EKS_INVENTORY_TIMEOUT is configured to 200ms, two 150ms calls exceed 200ms and must fail
	t.Setenv(KS_EKS_INVENTORY_TIMEOUT_ENV_VAR, "200ms")
	_, err := eksSupport.GetListEntitiesForPolicies("us-east-1")
	assert.Error(t, err)

	// When KS_EKS_INVENTORY_TIMEOUT is unset, individual calls succeed and the full inventory completes
	os.Unsetenv(KS_EKS_INVENTORY_TIMEOUT_ENV_VAR)
	out, err := eksSupport.GetListEntitiesForPolicies("us-east-1")
	assert.NoError(t, err)
	assert.Len(t, out.EntitiesForPolicies, 2)
}

func TestGetListEntitiesForPolicies_RegressionBeyond30s(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 30s+ regression test in short mode")
	}

	const policyCount = 62
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		action := r.Form.Get("Action")
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		switch action {
		case "ListPolicies":
			var members strings.Builder
			for i := 0; i < policyCount; i++ {
				members.WriteString(fmt.Sprintf(`
      <member>
        <Arn>arn:aws:iam::123456789012:policy/Policy-%d</Arn>
        <PolicyName>Policy-%d</PolicyName>
        <DefaultVersionId>v1</DefaultVersionId>
      </member>`, i, i))
			}
			w.Write([]byte(fmt.Sprintf(`
<ListPoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListPoliciesResult>
    <Policies>%s</Policies>
    <IsTruncated>false</IsTruncated>
  </ListPoliciesResult>
  <ResponseMetadata><RequestId>req-1</RequestId></ResponseMetadata>
</ListPoliciesResponse>`, members.String())))
		case "ListEntitiesForPolicy":
			// 500ms per entity response (62 * 500ms = 31 seconds)
			time.Sleep(500 * time.Millisecond)
			w.Write([]byte(`
<ListEntitiesForPolicyResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
  <ListEntitiesForPolicyResult>
    <PolicyGroups/><PolicyUsers/><PolicyRoles/><IsTruncated>false</IsTruncated>
  </ListEntitiesForPolicyResult>
  <ResponseMetadata><RequestId>req-2</RequestId></ResponseMetadata>
</ListEntitiesForPolicyResponse>`))
		}
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")
	os.Unsetenv(KS_EKS_INVENTORY_TIMEOUT_ENV_VAR)

	eksSupport := NewEKSSupport()
	start := time.Now()
	out, err := eksSupport.GetListEntitiesForPolicies("us-east-1")
	elapsed := time.Since(start)

	require.NoError(t, err, "enumeration beyond 30s must succeed without being cancelled by an artificial overall ceiling")
	require.NotNil(t, out)
	assert.Len(t, out.EntitiesForPolicies, policyCount)
	assert.Greater(t, elapsed, 30*time.Second, "total elapsed time must exceed 30 seconds to validate regression fix")
}
