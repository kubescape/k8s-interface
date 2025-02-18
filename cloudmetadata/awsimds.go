package cloudmetadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apitypes "github.com/armosec/armoapi-go/armotypes"
)

const (
	metadataEndpoint = "http://169.254.169.254"
	defaultTimeout   = 2 * time.Second
	tokenPath        = "/latest/api/token"
	tokenTTL         = "21600" // 6 hours in seconds
)

type MetadataClient struct {
	client    *http.Client
	useIMDSv2 bool
}

// NewMetadataClient creates a new client for fetching EC2 metadata
func NewMetadataClient(useIMDSv2 bool) *MetadataClient {
	return &MetadataClient{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		useIMDSv2: useIMDSv2,
	}
}

// getToken fetches an IMDSv2 token
func (m *MetadataClient) getToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "PUT", metadataEndpoint+tokenPath, nil)
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", tokenTTL)

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading token response: %w", err)
	}

	return string(token), nil
}

// get makes a GET request to the metadata service
func (m *MetadataClient) get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", metadataEndpoint+path, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if m.useIMDSv2 {
		token, err := m.getToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting IMDSv2 token: %w", err)
		}
		req.Header.Set("X-aws-ec2-metadata-token", token)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return data, nil
}

// getMetadataValue is a helper function to get a single metadata value
func (m *MetadataClient) getMetadataValue(ctx context.Context, path string) (string, error) {
	data, err := m.get(ctx, "/latest/meta-data/"+path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetMetadata fetches all available instance metadata
func (m *MetadataClient) GetMetadata(ctx context.Context) (*apitypes.CloudMetadata, error) {
	// Get instance identity document first
	data, err := m.get(ctx, "/latest/dynamic/instance-identity/document")
	if err != nil {
		return nil, fmt.Errorf("fetching instance identity document: %w", err)
	}

	var identityDoc struct {
		Region           string `json:"region"`
		InstanceID       string `json:"instanceId"`
		InstanceType     string `json:"instanceType"`
		AccountID        string `json:"accountId"`
		PrivateIP        string `json:"privateIp"`
		AvailabilityZone string `json:"availabilityZone"`
	}

	if err := json.Unmarshal(data, &identityDoc); err != nil {
		return nil, fmt.Errorf("parsing identity document: %w", err)
	}

	// Get public IP (might not be available)
	publicIP, _ := m.getMetadataValue(ctx, "public-ipv4")

	// Get hostname
	hostname, err := m.getMetadataValue(ctx, "hostname")
	if err != nil {
		// Fallback to local-hostname
		hostname, _ = m.getMetadataValue(ctx, "local-hostname")
	}

	metadata := &apitypes.CloudMetadata{
		Provider:     ProviderAWS,
		InstanceID:   identityDoc.InstanceID,
		InstanceType: identityDoc.InstanceType,
		Region:       identityDoc.Region,
		Zone:         identityDoc.AvailabilityZone,
		PrivateIP:    identityDoc.PrivateIP,
		PublicIP:     publicIP,
		Hostname:     hostname,
		AccountID:    identityDoc.AccountID,
	}

	return metadata, nil
}
