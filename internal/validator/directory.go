package validator

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ValidateDataDirectory(dir string) error {
	absoluteDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if absoluteDir == workDir {
		return errors.New("cannot use project root directory as data directory")
	}

	if strings.Contains(absoluteDir, filepath.Join(workDir, "internal")) {
		return errors.New("cannot use internal as data directory")
	}

	_, err = os.Stat(absoluteDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return validateExistingDir(absoluteDir)
}

func validateExistingDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	hasBucketsCSV := false

	for _, entry := range entries {
		name := entry.Name()
		if name == "buckets.csv" {
			hasBucketsCSV = true
			break
		}
	}

	if !hasBucketsCSV {
		return errors.New("directory can not be used as data directory")
	}

	return nil
}
