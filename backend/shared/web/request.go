package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func GetIntParam(query url.Values, key string, defaultValue int) int {
	if value := query.Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func Decode[T any](r *http.Request, v *T) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode request body: %w", err)
	}
	err := r.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
