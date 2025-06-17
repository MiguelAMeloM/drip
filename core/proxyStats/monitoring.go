/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package proxyStats

import (
	"fmt"
	"sync"
	"time"
)

// ANSI color codes
const (
	Reset    = "\033[0m"
	Bold     = "\033[1m"
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Cyan     = "\033[36m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

type ProxyStats struct {
	NumberOfRequests int         `json:"number_of_requests"`
	ResponseTime     float64     `json:"response_time"`
	StartedAt        time.Time   `json:"started_at"`
	ActiveServers    int         `json:"active_servers"`
	Mutex            *sync.Mutex `json:"-"`
}

func New() *ProxyStats {
	return &ProxyStats{Mutex: &sync.Mutex{}, StartedAt: time.Now()}
}

func (s *ProxyStats) Increment(start time.Time) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ResponseTime += time.Now().Sub(start).Seconds()
	s.NumberOfRequests++
}

func (s *ProxyStats) String() string {
	rps := float64(s.NumberOfRequests) / time.Since(s.StartedAt).Seconds()
	meanLatency := s.ResponseTime / float64(s.NumberOfRequests)

	// Color según latencia
	var latencyColor string
	switch {
	case meanLatency < 100:
		latencyColor = Green
	case meanLatency < 300:
		latencyColor = Yellow
	default:
		latencyColor = Red
	}

	// Color según RPS
	var rpsColor string
	switch {
	case rps > 1000:
		rpsColor = Green
	case rps > 100:
		rpsColor = Yellow
	default:
		rpsColor = Red
	}

	// Color según número de servidores
	var serverColor string
	switch {
	case s.ActiveServers >= 10:
		serverColor = Cyan
	case s.ActiveServers >= 5:
		serverColor = Yellow
	default:
		serverColor = Red
	}

	// Cuerpo formateado con colores
	body := fmt.Sprintf(
		"%s%15d%s  %s%15.2f%s  %s%15.2f%s  %s%15d%s",
		Reset, s.NumberOfRequests, Reset,
		latencyColor, meanLatency, Reset,
		rpsColor, rps, Reset,
		serverColor, s.ActiveServers, Reset,
	)

	return body
}

func (s *ProxyStats) Reset() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.StartedAt = time.Now()
	s.ResponseTime = 0
	s.NumberOfRequests = 0
}

func (s *ProxyStats) RequestsPerSecond() float64 {
	return float64(s.NumberOfRequests) / time.Since(s.StartedAt).Seconds()
}
