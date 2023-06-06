package names

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// instanceIDSlugHashlessFormat is a format of the Instance ID slug string without hash-based identifiers
	instanceIDSlugHashlessFormat = "%s-%s-%s"
	// instanceIDSlugFormat is a format of the slug string:
	// "hashlessFormat + hashLeading + hashTrailing"
	instanceIDSlugFormat     = "%s-%s-%s"
	instanceIDSlugHashLength = 4

	// imageIDSlugFormat is a format of the Image ID slug
	imageIDSlugFormat     = "%s-%s"
	imageIDSlugHashLength = 6

	maxDNSSubdomainLength = 253
	maxImageNameLength    = maxDNSSubdomainLength - imageIDSlugHashLength - 1
)

// imageToDNSSubdomainReplacer is a replacer that can replace a valid, well-formed container image string to a valid DNS Subdomain
var imageToDNSSubdomainReplacer = strings.NewReplacer("://", "-", ":", "-", "/", "-", "_", "-", "@", "-")

var (
	dnsSubdomainRegexp = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,251}[a-z0-9]$`)
	dnsLabelRegexp     = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
)

// IsValidDNSSubdomainName returns true if a given string is a valid DNS Subdomain name as defined in the Kubernetes docs
func IsValidDNSSubdomainName(s string) bool {
	return dnsSubdomainRegexp.MatchString(s)
}

// IsValidDNSLabelName returns true if a given string is a valid DNS label name as defined in the Kubernetes docs
func IsValidDNSLabelName(s string) bool {
	return dnsLabelRegexp.MatchString(s)
}

// IsValidSlug returns true if a given string is a valid slug, else false
//
// A string is considered a valid slug if it can be used as a name of a Kubernetes resource
func IsValidSlug(s string) bool {
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

// InstanceIDToSlug retuns a human-friendly representation given a description of an instance ID
//
// If the given inputs would produce an invalid slug, it returns an appropriate error
func InstanceIDToSlug(name, namespace, kind, hashedID string) (string, error) {
	leadingDigest, trailingDigest := hashedID[:instanceIDSlugHashLength], hashedID[len(hashedID)-instanceIDSlugHashLength:]

	hashlessInstanceIDSlug := fmt.Sprintf(instanceIDSlugHashlessFormat, namespace, kind, name)
	hashlessInstanceIDSlug = sanitizeInstanceIDSlug(hashlessInstanceIDSlug)

	var err error
	slug, err := fmt.Sprintf(instanceIDSlugFormat, hashlessInstanceIDSlug, leadingDigest, trailingDigest), nil
	slug = strings.ToLower(slug)

	if !IsValidSlug(slug) {
		slug, err = "", ErrInvalidSlug
	}

	return strings.ToLower(slug), err
}

// ImageInfoToSlug returns a human-friendly representation for a given image information
//
// If the given inputs would produce an invalid slug, it returns an appropriate error
func ImageInfoToSlug(image, imageHash string) (string, error) {
	if len(image) == 0 || len(imageHash) < imageIDSlugHashLength {
		return "", ErrInvalidSlug
	}

	var err error
	imageHashStub := imageHash[len(imageHash)-imageIDSlugHashLength:]
	sanitizedImage := sanitizeImage(image)
	slug, err := fmt.Sprintf(imageIDSlugFormat, sanitizedImage, imageHashStub), nil
	slug = strings.ToLower(slug)

	if !IsValidSlug(slug) {
		slug, err = "", ErrInvalidSlug
	}

	return slug, err
}
