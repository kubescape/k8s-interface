package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func isolateADC(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("CLOUDSDK_CONFIG", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)
	t.Setenv("APPDATA", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)
}

func TestGKETimeoutConstants(t *testing.T) {
	assert.Equal(t, 5*time.Second, gkeCallTimeout)
	assert.Equal(t, 30*time.Second, gkeRBACEnumerationTimeout)
	assert.Less(t, gkeCallTimeout, gkeRBACEnumerationTimeout)
	assert.Equal(t, 30*time.Second, tokenHTTPTimeout)
	// a token request gets at least as long as the longest caller waits for it
	assert.LessOrEqual(t, gkeRBACEnumerationTimeout, tokenHTTPTimeout)
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
	isolateADC(t)
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

	t.Setenv("GCE_METADATA_HOST", server.Listener.Addr().String())

	gkeSupport := NewGKESupport()
	start := time.Now()
	token, err := gkeSupport.GetAuthorizationKey()
	elapsed := time.Since(start)

	assert.Error(t, err)
	if token != "" {
		t.Fatalf("expected empty token on timeout, got token (value redacted for security, length: %d)", len(token))
	}
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	// Bounds check: must abort near 5s deadline and not escape to ~20s retry duration
	assert.Less(t, elapsed, 7*time.Second)
	assert.GreaterOrEqual(t, elapsed, 4500*time.Millisecond)
}

func TestGKEGetAuthorizationKey_Success(t *testing.T) {
	isolateADC(t)
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

	t.Setenv("GCE_METADATA_HOST", server.Listener.Addr().String())

	gkeSupport := NewGKESupport()
	token, err := gkeSupport.GetAuthorizationKey()
	assert.NoError(t, err)
	if token != "ya29.mock-test-access-token" {
		t.Fatalf("unexpected token returned (value redacted for security, length: %d)", len(token))
	}
}

func TestGKEGetAuthorizationKey_AuthorizedUserADC(t *testing.T) {
	tempDir := t.TempDir()
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": fmt.Sprintf("ya29.mock-token-%d", count),
			"token_type":   "Bearer",
			"expires_in":   1, // 1 second lifetime to trigger refresh on subsequent request
		})
	}))
	defer server.Close()

	adcJSON := fmt.Sprintf(`{
		"type": "authorized_user",
		"client_id": "mock-client-id",
		"client_secret": "mock-client-secret",
		"refresh_token": "mock-refresh-token",
		"token_uri": "%s"
	}`, server.URL)

	adcFile := filepath.Join(tempDir, "application_default_credentials.json")
	err := os.WriteFile(adcFile, []byte(adcJSON), 0600)
	assert.NoError(t, err)

	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", adcFile)
	t.Setenv("CLOUDSDK_CONFIG", tempDir)

	gkeSupport := NewGKESupport()

	// Initial token retrieval
	token1, err := gkeSupport.GetAuthorizationKey()
	assert.NoError(t, err)
	if token1 != "ya29.mock-token-1" {
		t.Fatalf("unexpected token returned (value redacted for security, length: %d)", len(token1))
	}
	assert.Equal(t, int32(1), atomic.LoadInt32(&requestCount))

	// Wait for token expiration (expires_in = 1s) to test refresh
	time.Sleep(1200 * time.Millisecond)

	token2, err := gkeSupport.GetAuthorizationKey()
	assert.NoError(t, err)
	if token2 != "ya29.mock-token-2" {
		t.Fatalf("unexpected refreshed token returned (value redacted for security, length: %d)", len(token2))
	}
	assert.Equal(t, int32(2), atomic.LoadInt32(&requestCount))
}

// A token endpoint that accepts the request and never answers must not leave
// requests open once callers give up. Each call still returns at the caller
// deadline; the token request behind it is cut off by tokenHTTPTimeout rather
// than staying open and piling up across calls.
func TestGKEGetAuthorizationKey_StalledTokenEndpointDrains(t *testing.T) {
	isolateADC(t)

	oldTimeout := tokenHTTPTimeout
	tokenHTTPTimeout = gkeCallTimeout + time.Second // outlives the caller, but not forever
	t.Cleanup(func() { tokenHTTPTimeout = oldTimeout })

	var inFlight, total int32
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The server only notices the client hanging up once the request
		// body is read, so read it before stalling.
		_, _ = io.Copy(io.Discard, r.Body)
		atomic.AddInt32(&total, 1)
		atomic.AddInt32(&inFlight, 1)
		defer atomic.AddInt32(&inFlight, -1)
		select {
		case <-r.Context().Done(): // the client gave up and closed the connection
		case <-release: // test cleanup
		}
	}))
	t.Cleanup(func() {
		close(release)
		server.Close()
	})

	adcJSON := fmt.Sprintf(`{
		"type": "authorized_user",
		"client_id": "mock-client-id",
		"client_secret": "mock-client-secret",
		"refresh_token": "mock-refresh-token",
		"token_uri": "%s"
	}`, server.URL)
	adcFile := filepath.Join(t.TempDir(), "application_default_credentials.json")
	assert.NoError(t, os.WriteFile(adcFile, []byte(adcJSON), 0600))
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", adcFile)

	gkeSupport := NewGKESupport()
	const calls = 2
	for i := 0; i < calls; i++ {
		start := time.Now()
		token, err := gkeSupport.GetAuthorizationKey()
		elapsed := time.Since(start)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		if token != "" {
			t.Fatalf("expected empty token on timeout, got token (value redacted for security, length: %d)", len(token))
		}
		assert.Less(t, elapsed, gkeCallTimeout+2*time.Second)
	}
	assert.Equal(t, int32(calls), atomic.LoadInt32(&total))

	assert.Eventually(t, func() bool { return atomic.LoadInt32(&inFlight) == 0 },
		tokenHTTPTimeout+2*time.Second, 50*time.Millisecond,
		"token requests still open after the callers gave up")
}

// The clients get a token source rather than loading ADC themselves, so the
// quota project from the ADC file has to be passed to them explicitly.
// GOOGLE_CLOUD_QUOTA_PROJECT keeps its precedence over the file.
func TestGKEGetListEntitiesForPolicies_QuotaProject(t *testing.T) {
	tests := []struct {
		name      string
		fileQuota string
		envQuota  string
		want      string
	}{
		{name: "from the ADC file", fileQuota: "quota-from-adc", want: "quota-from-adc"},
		{name: "env wins over the file", fileQuota: "quota-from-adc", envQuota: "quota-from-env", want: "quota-from-env"},
		{name: "env only", envQuota: "quota-from-env", want: "quota-from-env"},
		{name: "none", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateADC(t)

			var gotQuota, gotAuth string
			var policyCalls int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/token" {
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"access_token": "ya29.mock-token",
						"token_type":   "Bearer",
						"expires_in":   3600,
					})
					return
				}
				if strings.HasSuffix(r.URL.Path, ":getIamPolicy") {
					atomic.AddInt32(&policyCalls, 1)
					gotQuota = r.Header.Get("X-Goog-User-Project")
					gotAuth = r.Header.Get("Authorization")
					_, _ = w.Write([]byte(`{"version":3,"bindings":[]}`))
					return
				}
				http.NotFound(w, r)
			}))
			defer server.Close()

			adc := map[string]string{
				"type":          "authorized_user",
				"client_id":     "mock-client-id",
				"client_secret": "mock-client-secret",
				"refresh_token": "mock-refresh-token",
				"token_uri":     server.URL + "/token",
			}
			if tt.fileQuota != "" {
				adc["quota_project_id"] = tt.fileQuota
			}
			adcJSON, err := json.Marshal(adc)
			assert.NoError(t, err)
			adcFile := filepath.Join(t.TempDir(), "application_default_credentials.json")
			assert.NoError(t, os.WriteFile(adcFile, adcJSON, 0600))
			t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", adcFile)
			t.Setenv(quotaProjectEnvVar, tt.envQuota)

			gkeSupport := &GKESupport{crmEndpoint: server.URL + "/"}
			policy, err := gkeSupport.GetListEntitiesForPolicies("my-project")
			if !assert.NoError(t, err) {
				return
			}
			assert.NotNil(t, policy)
			assert.Equal(t, int32(1), atomic.LoadInt32(&policyCalls))
			assert.Equal(t, "Bearer ya29.mock-token", gotAuth, "the request must use the ADC token")
			assert.Equal(t, tt.want, gotQuota)
		})
	}
}

func TestGKEGetClusterDescribe_Timeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer listener.Close()

	gkeSupport := &GKESupport{
		containerEndpoint: listener.Addr().String(),
		tokenSource:       oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "mock-test-token"}),
	}

	start := time.Now()
	cluster, err := gkeSupport.GetClusterDescribe("test-cluster", "us-central1", "test-project")
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Nil(t, cluster)
	assert.Less(t, elapsed, 7*time.Second)
}

func TestGKEGetClusterDescribe_MetadataTimeout(t *testing.T) {
	testCases := []struct {
		name              string
		disableNewAuthLib bool
	}{
		{
			name:              "Default auth library",
			disableNewAuthLib: false,
		},
		{
			name:              "Legacy auth library (DISABLE_NEW_AUTH_LIB=true)",
			disableNewAuthLib: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateADC(t)
			if tc.disableNewAuthLib {
				t.Setenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_NEW_AUTH_LIB", "true")
			} else {
				t.Setenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_NEW_AUTH_LIB", "false")
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "project-id") {
					w.Header().Set("Metadata-Flavor", "Google")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("mock-project"))
					return
				}
				time.Sleep(3 * time.Second)
				w.Header().Set("Metadata-Flavor", "Google")
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			t.Setenv("GCE_METADATA_HOST", server.Listener.Addr().String())

			// Authentication is kept enabled with no mock tokenSource or endpoint bypass
			gkeSupport := NewGKESupport()
			start := time.Now()
			cluster, err := gkeSupport.GetClusterDescribe("test-cluster", "us-central1", "test-project")
			elapsed := time.Since(start)

			assert.Error(t, err)
			assert.Nil(t, cluster)
			assert.Contains(t, err.Error(), "context deadline exceeded")
			assert.Less(t, elapsed, 7*time.Second)
			assert.GreaterOrEqual(t, elapsed, 4500*time.Millisecond)
		})
	}
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
