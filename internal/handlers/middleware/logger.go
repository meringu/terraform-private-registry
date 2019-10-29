package middleware

import (
	"net/http"

	"github.com/google/go-github/github"
	tcontext "github.com/meringu/terraform-private-registry/internal/context"
)

// Used to capture the status code
type statusWriter struct {
	http.ResponseWriter
	status    int
	bytesSent int
}

// Header proxies the function to the encapsulated http.ResponseWriter
func (w *statusWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader proxies the function to the encapsulated http.ResponseWriter and captures the status
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write proxies the function to the encapsulated http.ResponseWriter and sets the status to 200 if not yet set
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	bytesSent, err := w.ResponseWriter.Write(b)
	w.bytesSent = bytesSent
	return bytesSent, err
}

// LoggerMiddleware logs the request
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := tcontext.GetLogger(r.Context())

		if r.Referer() != "" {
			logger = logger.WithField("http_referer", r.Referer())
		}
		if r.UserAgent() != "" {
			logger = logger.WithField("http_user_agent", r.UserAgent())
		}

		logger = logger.
			WithField("remote_addr", r.RemoteAddr).
			WithField("request", r.URL.RequestURI())

		if id := github.DeliveryID(r); id != "" {
			logger = logger.WithField("github_delivery_id", id)
		}

		s := &statusWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(s, r.WithContext(tcontext.WithLogger(r.Context(), logger)))

		logger = logger.
			WithField("body_bytes_sent", s.bytesSent).
			WithField("status", s.status)

		logger.Infof("response completed")
	})
}
