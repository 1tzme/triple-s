# Triple-S (Simple Storage Service)

A lightweight S3-compatible object storage server written in Go.

## Features

- Bucket management (create/list/delete)
- Object operations (upload/download/delete)
- S3-compatible XML API responses
- Local file system storage with CSV metadata

## Installation

```bash
go build -o triple-s .
```

## Usage

```bash
# Start server (default port 8080, data dir ./data)
./triple-s

# Custom port and data directory
./triple-s -port 7777 -dir ./storage
```

## API Examples

### Bucket Operations

```bash
# Create bucket
curl -X PUT http://localhost:8080/my-bucket

# List buckets 
curl http://localhost:8080/

# Delete bucket
curl -X DELETE http://localhost:8080/my-bucket
```

### Object Operations

```bash
# Upload file
curl -X PUT -T image.jpg http://localhost:8080/my-bucket/photo.jpg

# Download file
curl http://localhost:8080/my-bucket/photo.jpg -o photo.jpg

# Delete file
curl -X DELETE http://localhost:8080/my-bucket/photo.jpg
```

## Bucket Naming Rules

- 3-63 characters
- Lowercase letters, numbers, hyphens, dots
- No consecutive special characters

## Data Storage Structure
```
.
├── bucket1
│   ├── file1.txt
│   └── objects.csv
├── bucket2
│   ├── image.jpg
│   └── objects.csv
└── buckets.csv
```

## Help
```bash
./triple-s --help
```

# Made with ❤️ by [zaaripzha](https://platform.alem.school/git/zaaripzha) aka [1tzme](https://github.com/1tzme)

