package handlers

import (
	"triple-s/internal/structure"
)

type Handler struct {
	server *structure.Server
}

func NewHandler(server *structure.Server) *Handler {
	return &Handler{}
}
