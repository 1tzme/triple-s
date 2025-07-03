package structure

import (
	"encoding/xml"
	"time"
)

type Server struct {
	Dir  string
	Port string
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type Bucket struct {
	Name         string    `xml:"Name"`
	CreationTime time.Time `xml:"CreationTime"`
	LastModified time.Time `xml:"LastModified"`
	Status       string    `xml:"Status"`
}

type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

type ListAllBuckets struct {
	XMLName xml.Name `xml:"ListAllBuckets"`
	Owner   Owner    `xml:"Owner"`
	Buckets Buckets  `xml:"Buckets"`
}

type Object struct {
	ObjectKey    string    `xml:"ObjectKey"`
	Size         int64     `xml:"Size"`
	ContentType  string    `xml:"ContentType"`
	LastModified time.Time `xml:"LastModified"`
}

type Error struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}
