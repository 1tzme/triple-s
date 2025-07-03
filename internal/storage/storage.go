package storage

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"

	"triple-s/internal/structure"
)

const (
	bucketsCSV = "buckets.csv"
	objectsCSV = "objects.csv"
)

func CreateBucket(dataDir, bucketName string) error {
	bucketDir := filepath.Join(dataDir, bucketName)
	err := os.MkdirAll(bucketDir, 0o755)
	if err != nil {
		return err
	}

	bucket := structure.Bucket{
		Name:         bucketName,
		CreationTime: time.Now(),
		LastModified: time.Now(),
		Status:       "active",
	}

	return addBucketToCSV(dataDir, bucket)
}

func addBucketToCSV(dataDir string, bucket structure.Bucket) error {
	csvPath := filepath.Join(dataDir, bucketsCSV)

	needHeader := false
	_, err := os.Stat(csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			needHeader = true
		} else {
			return err
		}
	}

	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if needHeader {
		writer.Write([]string{"Name", "CreationTime", "LastModifiedTime", "Status"})
	}

	record := []string{
		bucket.Name,
		bucket.CreationTime.Format(time.RFC3339),
		bucket.LastModified.Format(time.RFC3339),
		bucket.Status,
	}
	return writer.Write(record)
}
