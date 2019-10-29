package handlers

import (
	"net/http"
)

const (
	pingMessage = "PONG\n"
)

// Ping returns PONG
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(pingMessage))
}

// PingHandler handles ping
func PingHandler() http.Handler {
	return http.HandlerFunc(Ping)
}
