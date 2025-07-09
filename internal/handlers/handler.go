package handlers

import (
	"encoding/xml"
	"net/http"

	"triple-s/internal/structure"
)

type Handler struct {
	server *structure.Server
}

func NewHandler(server *structure.Server) *Handler {
	return &Handler{
		server: server,
	}
}

func (h *Handler) sendError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)

	errorResp := structure.Error{
		Code:    code,
		Message: message,
	}

	xml.NewEncoder(w).Encode(errorResp)
}
