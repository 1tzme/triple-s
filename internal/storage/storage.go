package storage

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func BucketExists(dataDir, bucketName string) (bool, error) {
	buckets, err := ListBuckets(dataDir)
	if err != nil {
		return false, err
	}

	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			return true, nil
		}
	}

	return false, nil
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

func DeleteBucket(dataDir string, bucketName string) error {
	bucketDir := filepath.Join(dataDir, bucketName)
	err := os.RemoveAll(bucketDir)
	if err != nil {
		return err
	}

	return removeBucketFromCSV(dataDir, bucketName)
}

func removeBucketFromCSV(dataDir, bucketName string) error {
	csvPath := filepath.Join(dataDir, bucketsCSV)

	buckets, err := ListBuckets(dataDir)
	if err != nil {
		return err
	}

	filteredBuckets := []structure.Bucket{}
	for _, bucket := range buckets {
		if bucket.Name != bucketName {
			filteredBuckets = append(filteredBuckets, bucket)
		}
	}

	file, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Name", "CreationTime", "LastModifiedTime", "Status"})

	for _, bucket := range filteredBuckets {
		record := []string{
			bucket.Name,
			bucket.CreationTime.Format(time.RFC3339),
			bucket.LastModified.Format(time.RFC3339),
			bucket.Status,
		}
		writer.Write(record)
	}

	return nil
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

func StoreObject(dataDir, bucketName, objectKey string, data []byte, object structure.Object) error {
	objectPath := filepath.Join(dataDir, bucketName, objectKey)

	err := os.MkdirAll(filepath.Dir(objectPath), 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile(objectPath, data, 0o644)
	if err != nil {
		return err
	}

	return addObjectToCSV(dataDir, bucketName, object)
}

func addObjectToCSV(dataDir, bucketName string, object structure.Object) error {
	csvPath := filepath.Join(dataDir, bucketName, objectsCSV)

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
		writer.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})
	}

	record := []string{
		object.ObjectKey,
		strconv.FormatInt(object.Size, 10),
		object.ContentType,
		object.LastModified.Format(time.RFC3339),
	}

	return writer.Write(record)
}

func ObjectExists(dataDir, bucketName, objectKey string) (bool, error) {
	objectPath := filepath.Join(dataDir, bucketName, objectKey)
	_, err := os.Stat(objectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, err
}

func GetObject(dataDir, bucketName, objectKey string) ([]byte, error) {
	objectPath := filepath.Join(dataDir, bucketName, objectKey)
	return os.ReadFile(objectPath)
}

func GetObjectMetadata(dataDir, bucketName, objectKey string) (*structure.Object, error) {
	objects, err := listObjects(dataDir, bucketName)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		if object.ObjectKey == objectKey {
			return &object, nil
		}
	}

	return nil, errors.New("object not found")
}

func listObjects(dataDir, bucketName string) ([]structure.Object, error) {
	csvPath := filepath.Join(dataDir, bucketName, objectsCSV)

	_, err := os.Stat(csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []structure.Object{}, nil
		} else {
			return nil, err
		}
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

	objects := []structure.Object{}
	for i, record := range records {
		if i == 0 && len(record) > 0 && record[0] == "ObjectKey" {
			continue
		}

		if len(record) < 4 {
			log.Printf("Not enough fields in line %d: expected 4, got %d", i+1, len(record))
			continue
		}

		size, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Printf("Failed to parse Size in line %d: %v", i+1, err)
			continue
		}
		lastModified, err := time.Parse(time.RFC3339, record[3])
		if err != nil {
			log.Printf("Failed to parse LastModified in line %d: %v", i+1, err)
			continue
		}

		object := structure.Object{
			ObjectKey:    record[0],
			Size:         size,
			ContentType:  record[2],
			LastModified: lastModified,
		}
		objects = append(objects, object)
	}

	return objects, nil
}
