package validator

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateBucketName(name string) error {
	bucketNameRegex := regexp.MustCompile(`^[a-z0-9.-]+$`)

	if len(name) < 3 || len(name) > 63 {
		return errors.New("bucket name must be between 3 and 63 characters long")
	}
	if !bucketNameRegex.MatchString(name) {
		return errors.New("bucket name can contain only lowercase letters, numbers, hyphens, periods")
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return errors.New("bucket name can not start or end with hyphens")
	}
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return errors.New("bucket name can not start or end with periods")
	}
	if strings.Contains(name, "..") || strings.Contains(name, "--") {
		return errors.New("bucket name can not contain double periods or hyphens")
	}

	return nil
}
