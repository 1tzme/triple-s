package router

import (
	"net/http"

	h "triple-s/internal/handlers"
	s "triple-s/internal/structure"
)

func Router(server *s.Server) *http.ServeMux {
	mux := http.NewServeMux()
	h := h.NewHandler(server)

	mux.HandleFunc("PUT /{bucketName}", h.PutBucket)
	mux.HandleFunc("GET /{$}", h.GetBuckets)
	mux.HandleFunc("DELETE /{bucketName}", h.DeleteBucket)
	mux.HandleFunc("PUT /{bucketName}/{objectKey}", h.PutObject)
	mux.HandleFunc("GET /{bucketName}/{objectKey}", h.GetObject)
	mux.HandleFunc("DELETE /{bucketName}/{objectKey}", h.DeleteObject)

	return mux
}
