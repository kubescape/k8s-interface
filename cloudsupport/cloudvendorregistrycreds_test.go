package cloudsupport

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/assert"
)

func TestParseECRImageTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		wantID     string
		wantRegion string
		wantErr    bool
	}{
		{name: "valid", tag: "015253967648.dkr.ecr.eu-central-1.amazonaws.com/armo:1", wantID: "015253967648", wantRegion: "eu-central-1"},
		{name: "too few segments", tag: "foo.dkr.ecr", wantErr: true},
		{name: "empty", tag: "", wantErr: true},
		{name: "empty registry id", tag: ".dkr.ecr.eu-central-1.amazonaws.com/armo:1", wantErr: true},
		{name: "empty region", tag: "123.dkr.ecr..amazonaws.com/armo:1", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, region, err := parseECRImageTag(tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantID, id)
			assert.Equal(t, tt.wantRegion, region)
		})
	}
}

func TestParseECRAuthorizationToken(t *testing.T) {
	enc := func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }

	user, pass, err := parseECRAuthorizationToken(enc("AWS:secret"))
	assert.NoError(t, err)
	assert.Equal(t, "AWS", user)
	assert.Equal(t, "secret", pass)

	// only the first ':' is a delimiter
	user, pass, err = parseECRAuthorizationToken(enc("AWS:a:b"))
	assert.NoError(t, err)
	assert.Equal(t, "AWS", user)
	assert.Equal(t, "a:b", pass)

	_, _, err = parseECRAuthorizationToken(enc("no-delimiter"))
	assert.Error(t, err)

	_, _, err = parseECRAuthorizationToken("%%%not-base64%%%")
	assert.Error(t, err)
}

func TestGetLoginDetailsForECR_MalformedTagDoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		_, _, err := GetLoginDetailsForECR("foo.dkr.ecr")
		assert.Error(t, err)
	})
}

func TestExtractECRAuthorizationToken(t *testing.T) {
	tok := "dGVzdA=="

	t.Run("nil response", func(t *testing.T) {
		_, err := extractECRAuthorizationToken(nil, "015253967648")
		assert.Error(t, err)
	})

	t.Run("empty AuthorizationData", func(t *testing.T) {
		res := &ecr.GetAuthorizationTokenOutput{AuthorizationData: nil}
		_, err := extractECRAuthorizationToken(res, "015253967648")
		assert.Error(t, err)
	})

	t.Run("nil first element", func(t *testing.T) {
		res := &ecr.GetAuthorizationTokenOutput{AuthorizationData: []*ecr.AuthorizationData{nil}}
		_, err := extractECRAuthorizationToken(res, "015253967648")
		assert.Error(t, err)
	})

	t.Run("nil AuthorizationToken", func(t *testing.T) {
		res := &ecr.GetAuthorizationTokenOutput{AuthorizationData: []*ecr.AuthorizationData{{AuthorizationToken: nil}}}
		_, err := extractECRAuthorizationToken(res, "015253967648")
		assert.Error(t, err)
	})

	t.Run("valid token returned", func(t *testing.T) {
		res := &ecr.GetAuthorizationTokenOutput{AuthorizationData: []*ecr.AuthorizationData{{AuthorizationToken: &tok}}}
		got, err := extractECRAuthorizationToken(res, "015253967648")
		assert.NoError(t, err)
		assert.Equal(t, tok, got)
	})
}

func TestBuildACRRefreshTokenError_DoesNotLeakToken(t *testing.T) {
	secretToken := "super-secret-azure-aad-access-token"
	err := buildACRRefreshTokenError(errors.New("boom"), "myregistry.azurecr.io", "tenant-123")

	msg := err.Error()
	assert.NotContains(t, msg, secretToken)
	assert.Contains(t, msg, "myregistry.azurecr.io")
	assert.Contains(t, msg, "tenant-123")
	assert.Contains(t, msg, "boom")
}
