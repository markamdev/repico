package server

import (
	"errors"
	"net/http"
	"strconv"
)

func Create(port int, handler http.Handler) (http.Server, error) {
	if port < 1 {
		return http.Server{}, errors.New("invalid listening port")
	}

	return http.Server{Addr: ":" + strconv.Itoa(port), Handler: handler}, nil
}
