package validator

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("bucket name must be between 3 and 63 characters long")
	}

	bucketNameRegex := regexp.MustCompile(`^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$`)
	if !bucketNameRegex.MatchString(name) {
		return errors.New("bucket name can only contain lowercase letters, numbers, hyphens, periods and must start and end with alphanumeric char")
	}
	if strings.Contains(name, "..") || strings.Contains(name, "--") || strings.Contains(name, "-.") || strings.Contains(name, "-.") {
		return errors.New("bucket name cannot contain consecutive special characters")
	}

	return nil
}
