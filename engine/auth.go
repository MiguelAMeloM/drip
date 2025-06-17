package engine

import (
	"net"
	"net/http"
)

var (
	AuthToken string
)

func isAllowed(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}

	return host == "127.0.0.1" || host == "::1" || host == "localhost" || verifyToken(r)
}

func verifyToken(r *http.Request) bool {
	return r.Header.Get("X-Auth-Token") == AuthToken
}
