package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// WriteJSON writes an interface as JSON onto the response writer
func WriteJSON(w http.ResponseWriter, status int, i interface{}) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write([]byte(data))
	return err
}
