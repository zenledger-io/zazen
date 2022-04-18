package httputils

import (
	"encoding/json"
	"net/http"
)

var (
	JSONContentType = "application/json; charset=utf-8"
)

func SendJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", JSONContentType)
	return json.NewEncoder(w).Encode(v)
}
