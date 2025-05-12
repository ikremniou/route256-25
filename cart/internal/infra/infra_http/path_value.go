package infra_http

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetInt64PathValueGt0(r *http.Request, key string) (int64, error) {
	var rawValue = r.PathValue(key)
	pathValueInt64, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid path value: %s", key)
	}

	if pathValueInt64 <= 0 {
		return 0, fmt.Errorf("path value: %s should be greater than 0", key)
	}

	return pathValueInt64, nil
}
