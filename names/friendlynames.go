package names

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// friendlyNameInstanceIDHashlessFormat is a format of the friendly name string without hash-based identifiers
	friendlyNameInstanceIDHashlessFormat = "%s-%s-%s"
	// friendlyNameInstanceIDFormat is a format of the friendly name string:
	// "hashlessFormat + hashLeading + hashTrailing"
	friendlyNameInstanceIDFormat     = "%s-%s-%s"
	friendlyNameInstanceIDHashLength = 4

	// friendlyNameImageIDHashLength is the length of the hash suffix to append to the friendly name
	friendlyNameImageIDFormat     = "%s-%s"
	friendlyNameImageIDHashLength = 6

	maxDNSSubdomainLength = 253
	maxImageNameLength    = maxDNSSubdomainLength - friendlyNameImageIDHashLength - 1
)

// imageToDNSSubdomainReplacer is a replacer that can replace a valid, well-formed container image string to a valid DNS Subdomain
var imageToDNSSubdomainReplacer = strings.NewReplacer("://", "-", ":", "-", "/", "-")

var (
	dnsSubdomainRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.-]{0,251}[a-zA-Z0-9]$`)
	dnsLabelRegexp     = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$`)
)

// IsValidDNSSubdomainName returns true if a given string is a valid DNS Subdomain name as defined in the Kubernetes docs
func IsValidDNSSubdomainName(s string) bool {
	return dnsSubdomainRegexp.MatchString(s)
}

// IsValidDNSLabelName returns true if a given string is a valid DNS label name as defined in the Kubernetes docs
func IsValidDNSLabelName(s string) bool {
	return dnsLabelRegexp.MatchString(s)
}

// IsValidFriendlyName returns true if a given string is a valid friendly name, else false
//
// A string is considered a valid friendly name if it can be used as a name of a Kubernetes resource
func IsValidFriendlyName(s string) bool {
	return IsValidDNSSubdomainName(s)
}

// sanitizeImage returns a sanitized image string safe for use with K8s names
//
// It expects a valid image name string
func sanitizeImage(image string) string {
	sanitized := imageToDNSSubdomainReplacer.Replace(image)

	if len(image) >= maxImageNameLength {
		sanitized = sanitized[:maxImageNameLength]
	}
	return sanitized
}

// sanitizeInstanceIDSlug returns a sanitized instance ID slug
func sanitizeInstanceIDSlug(instanceIDSlug string) string {
	if len(instanceIDSlug) > 243 {
		return instanceIDSlug[:243]
	} else {
		return instanceIDSlug
	}

}

// InstanceIDToFriendlyName retuns a human-friendly name representation given a description of an instance ID
//
// If the given inputs would produce an invalid friendly name, it returns an appropriate error
func InstanceIDToFriendlyName(name, namespace, kind, hashedID string) (string, error) {
	leadingDigest, trailingDigest := hashedID[:friendlyNameInstanceIDHashLength], hashedID[len(hashedID)-friendlyNameInstanceIDHashLength:]

	hashlessInstanceIDSlug := fmt.Sprintf(friendlyNameInstanceIDHashlessFormat, namespace, kind, name)
	hashlessInstanceIDSlug = sanitizeInstanceIDSlug(hashlessInstanceIDSlug)

	var err error
	friendlyName, err := fmt.Sprintf(friendlyNameInstanceIDFormat, hashlessInstanceIDSlug, leadingDigest, trailingDigest), nil

	if !IsValidFriendlyName(friendlyName) {
		friendlyName, err = "", ErrInvalidFriendlyName
	}

	return friendlyName, err
}

// ImageInfoToFriendlyName returns a friendly name for a given image information
//
// If the given inputs would produce an invalid friendly name, it returns an appropriate error
func ImageInfoToFriendlyName(image, imageHash string) (string, error) {
	if len(image) == 0 || len(imageHash) < friendlyNameImageIDHashLength {
		return "", ErrInvalidFriendlyName
	}

	var err error
	imageHashStub := imageHash[len(imageHash)-friendlyNameImageIDHashLength:]
	sanitizedImage := sanitizeImage(image)
	friendlyName, err := fmt.Sprintf(friendlyNameImageIDFormat, sanitizedImage, imageHashStub), nil

	if !IsValidFriendlyName(friendlyName) {
		friendlyName, err = "", ErrInvalidFriendlyName
	}

	return friendlyName, err
}
