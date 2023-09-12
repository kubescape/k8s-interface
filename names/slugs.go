package names

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubescape/k8s-interface/workloadinterface"
)

const (
	// instanceIDSlugHashlessFormat is a format of the Instance ID slug string without hash-based identifiers
	instanceIDSlugHashlessFormat = "%s-%s-%s"
	// instanceIDSlugFormat is a format of the slug string:
	// "hashlessFormat + hashLeading + hashTrailing"
	slugFormat     = "%s-%s-%s"
	slugHashLength = 4
	// slugHashesLength is a length of the hash-based identifiers (-xxxx-xxxx) in the slug string
	slugHashesLength        = slugHashLength*2 + 2
	maxHashlessStringLength = maxDNSSubdomainLength - slugHashesLength

	// imageIDSlugFormat is a format of the Image ID slug
	imageIDSlugFormat     = "%s-%s"
	imageIDSlugHashLength = 6

	maxDNSSubdomainLength = 253
	maxImageNameLength    = maxDNSSubdomainLength - imageIDSlugHashLength - 1
)

// imageToDNSSubdomainReplacer is a replacer that can replace a valid, well-formed container image string to a valid DNS Subdomain
var imageToDNSSubdomainReplacer = strings.NewReplacer("://", "-", ":", "-", "/", "-", "_", "-", "@", "-")

var (
	dnsSubdomainRegexp         = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,251}[a-z0-9]$`)
	nonDnsSubdomainCharsRegexp = regexp.MustCompile(`[^a-zA-Z0-9\-\.]`)
	dnsLabelRegexp             = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
	labelValueRegexp           = regexp.MustCompile(`^$|^[a-zA-Z0-9]([-_.a-zA-Z0-9]{0,61}[a-zA-Z0-9])?$`)
	nonLabelValueCharsRegexp   = regexp.MustCompile(`[^a-zA-Z0-9\-_.]`)
)

func ToValidDNSSubdomainName(input string) (string, error) {
	if IsValidDNSSubdomainName(input) {
		return input, nil
	}

	// Replace non-DNS subdomain characters with hyphens.
	dnsCompatible := nonDnsSubdomainCharsRegexp.ReplaceAllString(input, "-")

	// Ensure that the result is in lowercase.
	dnsCompatible = strings.ToLower(dnsCompatible)

	// Trim leading and trailing hyphens.
	dnsCompatible = strings.Trim(dnsCompatible, "-")

	// Limit the length to 253 characters as required.
	if len(dnsCompatible) > 253 {
		dnsCompatible = dnsCompatible[:253]
	}

	// Ensure that the name starts and ends with an alphanumeric character.
	dnsCompatible = strings.TrimFunc(dnsCompatible, isNonAlphanumeric)
	if len(dnsCompatible) == 0 {
		return "", fmt.Errorf("cannot transform input into a valid DNS subdomain name: %s", input)
	}

	return dnsCompatible, nil
}

func ToValidLabelValue(input string) string {
	if IsValidLabelValue(input) {
		return input
	}

	// Replace non-label value characters with hyphens ("-").
	labelValueCompatible := nonLabelValueCharsRegexp.ReplaceAllString(input, "-")

	// Ensure that the result is in lowercase.
	labelValueCompatible = strings.ToLower(labelValueCompatible)

	// Trim leading and trailing hyphens.
	labelValueCompatible = strings.Trim(labelValueCompatible, "-")

	// Limit the length to 63 characters as required.
	if len(labelValueCompatible) > 63 {
		labelValueCompatible = labelValueCompatible[:63]
	}

	// Ensure that the label value begins and ends with an alphanumeric character.
	labelValueCompatible = strings.TrimFunc(labelValueCompatible, isNonAlphanumeric)

	return labelValueCompatible
}

// IsValidDNSSubdomainName returns true if a given string is a valid DNS Subdomain name as defined in the Kubernetes docs
func IsValidDNSSubdomainName(s string) bool {
	return dnsSubdomainRegexp.MatchString(s)
}

// IsValidDNSLabelName returns true if a given string is a valid DNS label name as defined in the Kubernetes docs
func IsValidDNSLabelName(s string) bool {
	return dnsLabelRegexp.MatchString(s)
}

// isValidLabelValue checks if a string is a valid Kubernetes label value as defined in the Kubernetes docs
func IsValidLabelValue(value string) bool {
	return labelValueRegexp.MatchString(value)
}

func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func isNonAlphanumeric(r rune) bool {
	return !isAlphanumeric(r)
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
	leadingDigest, trailingDigest := hashedID[:slugHashLength], hashedID[len(hashedID)-slugHashLength:]

	hashlessInstanceIDSlug := fmt.Sprintf(instanceIDSlugHashlessFormat, namespace, kind, name)
	hashlessInstanceIDSlug = sanitizeInstanceIDSlug(hashlessInstanceIDSlug)

	var err error
	slug, err := fmt.Sprintf(slugFormat, hashlessInstanceIDSlug, leadingDigest, trailingDigest), nil
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

func GetNamespaceLessSlug(slug, namespace string) (string, error) {
	if !IsValidSlug(slug) {
		return "", ErrInvalidSlug
	}
	namespaceLessSlug := strings.TrimPrefix(slug, namespace+"-")
	if slug == namespaceLessSlug {
		return "", ErrInvalidSlug
	}
	if !IsValidSlug(namespaceLessSlug) {
		return "", ErrInvalidSlug
	}

	return namespaceLessSlug, nil
}

func SanitizeLabelValues(labels map[string]string) {
	for k, v := range labels {
		labels[k] = ToValidLabelValue(v)
	}
}

// StringToSlug receives any string and returns a human-friendly representation of it as a slug
//
// If the given inputs would produce an invalid slug, it returns an appropriate error
func StringToSlug(str string) (string, error) {
	// hash the string, take the first and last 4 characters of the hash
	hashBytes := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hashBytes[:])
	leadingDigest, trailingDigest := hashStr[:slugHashLength], hashStr[len(hashStr)-slugHashLength:]

	// sanitize the string to be DNS Subdomain compatible
	sanitizedStr, err := ToValidDNSSubdomainName(str)
	if err != nil {
		return "", err
	}

	if len(sanitizedStr) >= maxHashlessStringLength {
		sanitizedStr = sanitizedStr[:maxHashlessStringLength]
	}

	slug := fmt.Sprintf(slugFormat, sanitizedStr, leadingDigest, trailingDigest)
	slug = strings.ToLower(slug)

	if !IsValidSlug(slug) {
		return "", ErrInvalidSlug
	}

	return strings.ToLower(slug), nil
}

func resourceToFormattedString(resource workloadinterface.IMetadata) string {
	return fmt.Sprintf("%s-%s-%s-%s", resource.GetApiVersion(), resource.GetKind(), resource.GetNamespace(), resource.GetName())
}

// ResourceToSlug returns a human-friendly representation for a given resource
// The slug is generated based on the API version, kind, namespace and name of the resource
func ResourceToSlug(resource workloadinterface.IMetadata) (string, error) {
	return StringToSlug(resourceToFormattedString(resource))
}

// RoleBindingResourceToSlug returns a human-friendly representation for a given role binding resource
// The slug is generated based on the subject, role and role binding resources of the role binding
func RoleBindingResourceToSlug(subject workloadinterface.IMetadata, role workloadinterface.IMetadata, roleBinding workloadinterface.IMetadata) (string, error) {
	subjectName := resourceToFormattedString(subject)
	roleName := resourceToFormattedString(role)
	roleBindingName := resourceToFormattedString(roleBinding)
	return StringToSlug(fmt.Sprintf("%s-%s-%s", subjectName, roleName, roleBindingName))
}
