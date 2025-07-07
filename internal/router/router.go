package router

import (
	"net/http"

	h "triple-s/internal/handlers"
	s "triple-s/internal/structure"
)

func Router(server *s.Server) *http.ServeMux {
	mux := http.NewServeMux()
	handler := h.NewHandler(server)

	mux.HandleFunc("PUT /{bucketName}", handler.PutBucket)
	mux.HandleFunc("GET /{$}", handler.GetBuckets)
	mux.HandleFunc("DELETE /{bucketName}", handler.DeleteBucket)
	mux.HandleFunc("PUT /{bucketName}/{objectKey}", handler.PutObject)
	mux.HandleFunc("GET /{bucketName}/{objectKey}", handler.GetObject)
	mux.HandleFunc("DELETE /{bucketName}/{objectKey}", handler.DeleteObject)

	return mux
}
