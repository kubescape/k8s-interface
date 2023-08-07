package names

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

func NormalizeImageName(image string) (string, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return image, fmt.Errorf("failed to parse image name %q: %w", image, err)
	} else {
		return named.String(), nil
	}
}
