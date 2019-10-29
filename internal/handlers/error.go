package handlers

import (
	"net/http"

	"github.com/meringu/terraform-private-registry/internal/api"
)

// NotFoundHandler handles not found
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.WriteNotFound(w)
	})
}

// NotImplementedHandler handles not implemented
func NotImplementedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.WriteNotImplemented(w)
	})
}
