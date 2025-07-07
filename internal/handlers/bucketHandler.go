package handlers

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/storage"
	"triple-s/internal/structure"
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

func (h *Handler) GetBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := storage.ListBuckets(h.server.Dir)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to list buckets", http.StatusInternalServerError)
		return
	}

	response := structure.ListAllBuckets{
		Owner: structure.Owner{
			ID:          "",
			DisplayName: "",
		},
		Buckets: structure.Buckets{
			Bucket: make([]structure.Bucket, len(buckets)),
		},
	}

	for _, bucket := range buckets {
		response.Buckets.Bucket = append(response.Buckets.Bucket, structure.Bucket{
			Name:         bucket.Name,
			CreationTime: bucket.CreationTime,
		})
	}

	w.Header().Set("Content Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	xml.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("bucketName")

	exists, err := storage.BucketExists(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "Internal error", "Failed to check bucket existence", http.StatusInternalServerError)
		return
	}
	if !exists {
		h.sendError(w, "NoSuchBucket", "The specified bucket does not exist", http.StatusNotFound)
		return
	}

	isEmpty, err := storage.IsBucketEmpty(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to check bucke emptiness", http.StatusInternalServerError)
		return
	}
	if !isEmpty {
		h.sendError(w, "BucketNotEmpty", "Bucket is not empty", http.StatusConflict)
		return
	}

	err = storage.DeleteBucket(h.server.Dir, bucketName)
	if err != nil {
		h.sendError(w, "InternalError", "Failed to delete bucket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
