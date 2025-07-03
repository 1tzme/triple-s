package storage

import (
	"encoding/csv"
	"log"
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

func ListBuckets(dataDir string) ([]structure.Bucket, error) {
	csvPath := filepath.Join(dataDir, bucketsCSV)
	_, err := os.Stat(csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []structure.Bucket{}, nil
		}
		return nil, err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	buckets := []structure.Bucket{}
	for i, record := range records {
		if i == 0 && len(record) > 0 && record[0] == "Name" {
			continue
		}
		if len(record) < 4 {
			log.Printf("Not enough fields in line %d: expected 4, got %d", i+1, len(record))
			continue
		}

		creationTime, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			log.Printf("Failed to parse CreationTime in line %d: %v", i+1, err)
			continue
		}
		modifiedTime, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			log.Printf("Failed to parse ModifiedTime in line %d: %v", i+1, err)
			continue
		}

		bucket := structure.Bucket{
			Name:         record[0],
			CreationTime: creationTime,
			LastModified: modifiedTime,
			Status:       record[3],
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
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
