/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package core

import (
	"bytes"
	"drip/core/modelProxy"
	"drip/core/proxyStats"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const (
	Stable    = "stable"
	ABtesting = "abtesting"
	Shadow    = "shadow"
	Canary    = "canary"
)

var (
	IntervalForMonitoring         = 60 * time.Second
	MaxRequestsPerSecond  float64 = 5000
	MinRequestsPerSecond  float64 = 100
)

func UpdateAutoscalingParams(interval time.Duration, minReq float64, maxReq float64) error {
	if interval < 60*time.Second {
		return errors.New("interval must be greater than 60s")
	}
	if maxReq < minReq {
		return errors.New("maxReq must be greater than minReq")
	}
	if maxReq < 0 || minReq < 0 {
		return errors.New("maxReq and minReq must be greater than 0")
	}
	IntervalForMonitoring = interval
	MinRequestsPerSecond = minReq
	MaxRequestsPerSecond = maxReq
	return nil
}

var Manager = &modelsManager{
	ModelsMapping: make(map[string]*modelProxy.LoadBalancer),
}

type modelsManager struct {
	ModelsMapping map[string]*modelProxy.LoadBalancer
}

func (mm *modelsManager) AddModel(ms modelProxy.NewModelSetting) error {
	switch ms.ReleaseType {
	case Stable:
		proxy := modelProxy.NewLoadBalancer(ms, modelProxy.NewStableRelease)
		mm.ModelsMapping[ms.ModelName] = proxy
		return nil
	case ABtesting:
		proxy := modelProxy.NewLoadBalancer(ms, modelProxy.NewABTestingProxy)
		mm.ModelsMapping[ms.ModelName] = proxy
		return nil
	case Shadow:
		proxy := modelProxy.NewLoadBalancer(ms, modelProxy.NewShadowProxy)
		mm.ModelsMapping[ms.ModelName] = proxy
		return nil
	case Canary:
		proxy := modelProxy.NewLoadBalancer(ms, modelProxy.NewCanaryProxy)
		mm.ModelsMapping[ms.ModelName] = proxy
		return nil
	default:
		return errors.New("model release type not supported")
	}
}

func (mm *modelsManager) RemoveModel(modelName string) error {
	if proxy, ok := mm.ModelsMapping[modelName]; ok {
		proxy.Close()
		delete(mm.ModelsMapping, modelName)
		return nil
	} else {
		return errors.New("model not found")
	}
}

func (mm *modelsManager) ForwardRequest(modelName string, req *http.Request) (gin.H, error) {
	if proxy, ok := mm.ModelsMapping[modelName]; ok {
		return proxy.ForwardRequest(req)
	} else {
		return nil, errors.New("model not found")
	}
}

func (mm *modelsManager) GenStats() (io.Reader, error) {
	stats := make(map[string]*proxyStats.ProxyStats)
	for modelName, proxy := range mm.ModelsMapping {
		stats[modelName] = proxy.GenStats()
	}

	b, err := json.Marshal(stats)
	if err != nil {
		return nil, err
	} else {
		return bytes.NewReader(b), nil
	}
}

func AutoScaling() {
	for {
		time.Sleep(IntervalForMonitoring)
		for _, cb := range Manager.ModelsMapping {
			stats := cb.GenStats().RequestsPerSecond()
			nOfServ := cb.Len()
			nOfServProxy := stats / float64(nOfServ)
			if nOfServProxy > MaxRequestsPerSecond {
				go cb.AddProxy()
			} else if nOfServProxy < MinRequestsPerSecond && nOfServ > 1 {
				go cb.RemoveProxy()
			}
		}
	}
}
