/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package modelProxy

import (
	"crypto/rand"
	"drip/core/serverRaiser"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os/exec"
	"syscall"
	"time"
)

const baseUrl string = "http://localhost"

type NewModelSetting struct {
	ModelName   string
	Raiser      serverRaiser.RaiseServer
	ReleaseType string
	Deadline    time.Time
	Prob        float64
}

type NewProxy func(NewModelSetting) Proxy

type Proxy interface {
	ForwardRequest(*http.Request) (gin.H, error)
	Close()
}

type ProxyBase struct {
	activeServers []*exec.Cmd
	modelName     string
	endpoint      string
}

func getAvailablePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func shutDownServer(p *ProxyBase) {
	for _, cmd := range p.activeServers {
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			cmd.Process.Kill()
		}
		time.Sleep(1 * time.Second)
	}
}

func genReqId() string {
	return rand.Text() + "-" + time.Now().Format("20060102150405")
}
