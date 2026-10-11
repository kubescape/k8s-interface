package v1

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGKETimeoutConstants(t *testing.T) {
	assert.Equal(t, 5*time.Second, gkeCallTimeout)
	assert.Equal(t, 30*time.Second, gkeRBACEnumerationTimeout)
	assert.Less(t, gkeCallTimeout, gkeRBACEnumerationTimeout)
}

func TestGKEGetRegion(t *testing.T) {
	gkeSupport := NewGKESupport()

	t.Run("Region extracted from valid cluster name", func(t *testing.T) {
		region, err := gkeSupport.GetRegion("gke_my-project_us-central1-a_my-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "us-central1-a", region)
	})

	t.Run("Region extracted from environment variable", func(t *testing.T) {
		t.Setenv(KS_CLOUD_REGION_ENV_VAR, "europe-west1")
		region, err := gkeSupport.GetRegion("gke_my-project_us-central1-a_my-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "europe-west1", region)
	})

	t.Run("Invalid cluster name returns error", func(t *testing.T) {
		if val, exists := os.LookupEnv(KS_CLOUD_REGION_ENV_VAR); exists {
			os.Unsetenv(KS_CLOUD_REGION_ENV_VAR)
			t.Cleanup(func() { os.Setenv(KS_CLOUD_REGION_ENV_VAR, val) })
		}
		_, err := gkeSupport.GetRegion("invalid-cluster-name")
		assert.Error(t, err)
	})
}

func TestGKEGetProject(t *testing.T) {
	gkeSupport := NewGKESupport()

	t.Run("Project extracted from valid cluster name", func(t *testing.T) {
		project, err := gkeSupport.GetProject("gke_my-project_us-central1-a_my-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "my-project", project)
	})

	t.Run("Project extracted from environment variable", func(t *testing.T) {
		t.Setenv(KS_GKE_PROJECT_ENV_VAR, "override-project")
		project, err := gkeSupport.GetProject("gke_my-project_us-central1-a_my-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "override-project", project)
	})

	t.Run("Invalid cluster name returns error", func(t *testing.T) {
		if val, exists := os.LookupEnv(KS_GKE_PROJECT_ENV_VAR); exists {
			os.Unsetenv(KS_GKE_PROJECT_ENV_VAR)
			t.Cleanup(func() { os.Setenv(KS_GKE_PROJECT_ENV_VAR, val) })
		}
		_, err := gkeSupport.GetProject("invalid")
		assert.Error(t, err)
	})
}

func TestGKEGetContextName(t *testing.T) {
	gkeSupport := NewGKESupport()

	t.Run("Context extracted from valid 4-part cluster name", func(t *testing.T) {
		contextName := gkeSupport.GetContextName("gke_my-project_us-central1-a_my-cluster")
		assert.Equal(t, "my-cluster", contextName)
	})

	t.Run("Invalid cluster name returns empty string", func(t *testing.T) {
		contextName := gkeSupport.GetContextName("invalid_cluster")
		assert.Equal(t, "", contextName)
	})
}

func TestGKEGetAuthorizationKey_MetadataTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "project-id") {
			w.Header().Set("Metadata-Flavor", "Google")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mock-project"))
			return
		}
		// Metadata token request delayed 3s then returning retryable HTTP 500
		time.Sleep(3 * time.Second)
		w.Header().Set("Metadata-Flavor", "Google")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("GCE_METADATA_HOST", server.Listener.Addr().String())

	gkeSupport := NewGKESupport()
	start := time.Now()
	token, err := gkeSupport.GetAuthorizationKey()
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	// Bounds check: must abort near 5s deadline and not escape to ~20s retry duration
	assert.Less(t, elapsed, 7*time.Second)
	assert.GreaterOrEqual(t, elapsed, 4500*time.Millisecond)
}

func TestGKEGetAuthorizationKey_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Metadata-Flavor", "Google")
		if strings.Contains(r.URL.Path, "project-id") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mock-project"))
			return
		}
		if strings.Contains(r.URL.Path, "token") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "ya29.mock-test-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("GCE_METADATA_HOST", server.Listener.Addr().String())

	gkeSupport := NewGKESupport()
	token, err := gkeSupport.GetAuthorizationKey()
	assert.NoError(t, err)
	assert.Equal(t, "ya29.mock-test-access-token", token)
}

func TestGKEGetClusterDescribe_Timeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer listener.Close()

	gkeSupport := &GKESupport{
		containerEndpoint: listener.Addr().String(),
	}

	start := time.Now()
	cluster, err := gkeSupport.GetClusterDescribe("test-cluster", "us-central1", "test-project")
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Nil(t, cluster)
	assert.Less(t, elapsed, 7*time.Second)
}

func TestGKEGetIAMMappings_Pagination(t *testing.T) {
	rolesPage1 := map[string]interface{}{
		"roles": []map[string]interface{}{
			{"name": "projects/test-project/roles/role1", "title": "Role One"},
		},
		"nextPageToken": "roles-token-2",
	}
	rolesPage2 := map[string]interface{}{
		"roles": []map[string]interface{}{
			{"name": "projects/test-project/roles/role2", "title": "Role Two"},
		},
	}
	saPage1 := map[string]interface{}{
		"accounts": []map[string]interface{}{
			{"name": "projects/test-project/serviceAccounts/sa1@test.iam.gserviceaccount.com"},
		},
		"nextPageToken": "sa-token-2",
	}
	saPage2 := map[string]interface{}{
		"accounts": []map[string]interface{}{
			{"name": "projects/test-project/serviceAccounts/sa2@test.iam.gserviceaccount.com"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		pageToken := r.URL.Query().Get("pageToken")

		if strings.Contains(path, "/roles") {
			if pageToken == "" {
				_ = json.NewEncoder(w).Encode(rolesPage1)
				return
			} else if pageToken == "roles-token-2" {
				_ = json.NewEncoder(w).Encode(rolesPage2)
				return
			}
		} else if strings.Contains(path, "/serviceAccounts") {
			if pageToken == "" {
				_ = json.NewEncoder(w).Encode(saPage1)
				return
			} else if pageToken == "sa-token-2" {
				_ = json.NewEncoder(w).Encode(saPage2)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	gkeSupport := &GKESupport{
		iamEndpoint: server.URL,
	}

	roleMappings, saMappings, err := gkeSupport.GetIAMMappings("test-project")
	assert.NoError(t, err)
	assert.Len(t, roleMappings, 2)
	assert.Equal(t, "Role One", roleMappings["projects/test-project/roles/role1"])
	assert.Equal(t, "Role Two", roleMappings["projects/test-project/roles/role2"])
	assert.Len(t, saMappings, 2)
	assert.Equal(t, "projects/test-project/serviceAccounts/sa1@test.iam.gserviceaccount.com", saMappings["projects/test-project/serviceAccounts/sa1@test.iam.gserviceaccount.com"])
	assert.Equal(t, "projects/test-project/serviceAccounts/sa2@test.iam.gserviceaccount.com", saMappings["projects/test-project/serviceAccounts/sa2@test.iam.gserviceaccount.com"])
}

func TestGKEGetIAMMappings_SubsequentPageError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		if strings.Contains(r.URL.Path, "/roles") {
			if pageToken == "" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"roles": []map[string]interface{}{
						{"name": "projects/test-project/roles/role1", "title": "Role One"},
					},
					"nextPageToken": "page2",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	gkeSupport := &GKESupport{
		iamEndpoint: server.URL,
	}

	roleMappings, saMappings, err := gkeSupport.GetIAMMappings("test-project")
	assert.Error(t, err)
	assert.Nil(t, roleMappings)
	assert.Nil(t, saMappings)
}
