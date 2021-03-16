package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteResponseObject(w http.ResponseWriter, r *http.Request, object interface{}) {
	if jsonBytes, err := json.Marshal(object); err != nil {
		return
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		_, _ = w.Write(jsonBytes)
	}
}
