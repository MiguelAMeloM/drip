/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

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
