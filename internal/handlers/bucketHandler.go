package handlers

import (
	"net/http"

	"triple-s/internal/storage"
	v "triple-s/internal/validator"
)

func (h *Handler) PutBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("bucketName")

	err := v.ValidateBucketName(bucketName)
	if err != nil {
		h.sendError(w, "InvalidBucketName", err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := storage.BucketExists(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to check bucket existence", http.StatusInternalServerError)
		return
	}
	if exists {
		h.sendError(w, "BucketAlreadyExists", "Bucket already exists", http.StatusConflict)
		return
	}

	err = storage.CreateBucket(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to create bucket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
