package router

import (
	"net/http"

	h "triple-s/internal/handlers"
	s "triple-s/internal/structure"
)

func Router(server *s.Server) *http.ServeMux {
	mux := http.NewServeMux()
	h := h.NewHandler(server)

	mux.HandleFunc("PUT /{BucketName}", h.PutBucket)
	mux.HandleFunc("GET /{$}", h.GetBuckets)
	mux.HandleFunc("DELETE /{BucketName}", h.DeleteBucket)
	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", h.PutObject)
	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", h.GetObject)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", h.DeleteObject)

	return mux
}
