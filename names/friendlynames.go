package names

import (
	"fmt"
	"strings"
)

const (
	// friendlyNameInstanceIDFormat is a format of the friendly name string
	friendlyNameInstanceIDFormat     = "%s-%s-%s-%s-%s"
	friendlyNameInstanceIDHashLength = 4

	// friendlyNameImageIDHashLength is the length of the hash suffix to append to the friendly name
	friendlyNameImageIDFormat     = "%s-%s"
	friendlyNameImageIDHashLength = 6
)

// imageToDNSSubdomainReplacer is a replacer that can replace a valid, well-formed container image string to a valid DNS Subdomain
var imageToDNSSubdomainReplacer = strings.NewReplacer("://", "-", ":", "-", "/", "-")

func isInvalidFriendlyName(s string) bool {
	return strings.ContainsAny(s, "/:")
}

// sanitizeImage returns a sanitized string safe for use with K8s names
func sanitizeImage(image string) string {
	return imageToDNSSubdomainReplacer.Replace(image)
}

// InstanceIDToFriendlyName retuns a human-friendly name representation given a description of an instance ID
//
// If the given inputs would produce an invalid friendly name, it returns an appropriate error
func InstanceIDToFriendlyName(name, namespace, kind, hashedID string) (string, error) {
	leadingDigest, trailingDigest := hashedID[:friendlyNameInstanceIDHashLength], hashedID[len(hashedID)-friendlyNameInstanceIDHashLength:]

	var err error
	friendlyName, err := fmt.Sprintf(friendlyNameInstanceIDFormat, namespace, kind, name, leadingDigest, trailingDigest), nil

	if isInvalidFriendlyName(friendlyName) {
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

	if isInvalidFriendlyName(friendlyName) {
		friendlyName, err = "", ErrInvalidFriendlyName
	}

	return friendlyName, err
}
