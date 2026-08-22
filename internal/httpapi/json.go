package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/rr173/task150-gridguard/internal/model"
	"net/http"
)

func decode(r *http.Request, target any) error {
	if r.Body == nil {
		return model.FieldError("body", "请求体不能为空")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
func write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func fail(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if model.IsNotFound(err) {
		status = http.StatusNotFound
	}
	var domain model.Error
	if errors.As(err, &domain) && domain.Field == "conflict" {
		status = http.StatusConflict
	}
	write(w, status, map[string]string{"error": err.Error()})
}
func method(w http.ResponseWriter, r *http.Request, want string) bool {
	if r.Method == want {
		return true
	}
	w.Header().Set("Allow", want)
	fail(w, model.FieldError("method", "需要 "+want))
	return false
}
