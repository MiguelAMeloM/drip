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
	"fmt"
	"github.com/MiguelAMeloM/drip/core/proxyStats"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var (
	serverRaisingDelay time.Duration = 5
)

func SetRaisingDelay(delay int) {
	serverRaisingDelay = time.Duration(delay)
}

type LoadBalancer struct {
	idx         int
	mutex       *sync.Mutex
	setting     NewModelSetting
	constructor NewProxy
	Proxies     []Proxy
	Stats       *proxyStats.ProxyStats
}

func NewLoadBalancer(setting NewModelSetting, constructor NewProxy) *LoadBalancer {
	newCB := &LoadBalancer{
		idx:         0,
		mutex:       &sync.Mutex{},
		setting:     setting,
		constructor: constructor,
		Stats:       proxyStats.New(),
	}
	newCB.AddProxy()
	return newCB
}

func (cb *LoadBalancer) GetUrls() []string {
	urls := make([]string, 0)
	for _, proxy := range cb.Proxies {
		urls = append(urls, proxy.GetUrls()...)
	}
	return urls
}

func (cb *LoadBalancer) Idx() int {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	ci := cb.idx
	if ci >= len(cb.Proxies) {
		ci = 0
	}

	if ci == len(cb.Proxies)-1 {
		cb.idx = 0
	} else {
		cb.idx++
	}

	return ci
}

func (cb *LoadBalancer) ForwardRequest(r *http.Request) (gin.H, error) {
	started := time.Now()
	defer cb.Stats.Increment(started)
	cb.mutex.Lock()
	idx := cb.Idx()
	proxy := cb.Proxies[idx]
	cb.mutex.Unlock()
	return proxy.ForwardRequest(r)
}

func (cb *LoadBalancer) Close() {
	for _, p := range cb.Proxies {
		p.Close()
	}
}

func (cb *LoadBalancer) GenStats() *proxyStats.ProxyStats {
	return cb.Stats
}

func (cb *LoadBalancer) Len() int {
	return len(cb.Proxies)
}

func (cb *LoadBalancer) AddProxy() {
	fmt.Printf("adding proxy %s\n", cb.setting.ModelName)
	proxy := cb.constructor(cb.setting)
	cb.Stats.Reset()
	time.Sleep(serverRaisingDelay * time.Second)
	cb.mutex.Lock()
	cb.Proxies = append(cb.Proxies, proxy)
	cb.Stats.ActiveServers++
	cb.mutex.Unlock()
}

func (cb *LoadBalancer) RemoveProxy() {
	fmt.Printf("removing proxy %d - %s\n", cb.idx, cb.setting.ModelName)
	i := len(cb.Proxies) - 1
	if i == 0 {
		return
	}
	proxy := cb.Proxies[i]
	cb.mutex.Lock()
	cb.Proxies = cb.Proxies[:i]
	cb.Stats.ActiveServers--
	cb.mutex.Unlock()
	proxy.Close()
	cb.Stats.Reset()
}
