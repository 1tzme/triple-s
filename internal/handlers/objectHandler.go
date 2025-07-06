package handlers

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"triple-s/internal/storage"
	"triple-s/internal/structure"
)

func (h *Handler) PutObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("buckeetName")
	objectKey := r.PathValue("objectKey")

	exists, err := storage.BucketExists(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to check bucket existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		h.sendError(w, "NoSuchBucket", "The specified bucket does not exists", http.StatusNotFound)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	contentLenStr := r.Header.Get("Content-Length")
	var contentLen int64
	if contentLenStr != "" {
		_, err := strconv.ParseInt(contentLenStr, 10, 64)
		if err != nil {
			h.sendError(w, "InvalidContentLength", "Content length is not valid", http.StatusBadRequest)
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if contentLen == 0 {
		contentLen = int64(len(body))
	}

	object := structure.Object{
		ObjectKey:    objectKey,
		Size:         contentLen,
		ContentType:  contentType,
		LastModified: time.Now(),
	}

	err = storage.StoreObject(h.server.Dir, bucketName, objectKey, body, object)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to store object", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
