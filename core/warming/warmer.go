/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package warming

import (
	"bytes"
	"encoding/json"
	"github.com/MiguelAMeloM/drip/core/modelProxy"
	"net/http"
	"sync"
	"time"
)

var (
	Manager warmersManager
	Delay   time.Duration = 24 * time.Hour
	mutex   *sync.Mutex   = &sync.Mutex{}
)

type warmersManager map[string]*Warmer

func (wm warmersManager) AddWarmer(modelName string, loadBalancer *modelProxy.LoadBalancer, delay time.Duration) {
	mutex.Lock()
	wm[modelName] = &Warmer{
		loadBalancer: loadBalancer,
		delay:        delay,
	}
	mutex.Unlock()
}

func (wm warmersManager) RemoveWarmer(modelName string) {
	mutex.Lock()
	delete(wm, modelName)
	mutex.Unlock()
}

type Warmer struct {
	loadBalancer *modelProxy.LoadBalancer
	delay        time.Duration
	sampleData   map[string]any //todo implementar lógica para catchear un sample
}

func (w *Warmer) WarmUp() {
	for _, url := range w.loadBalancer.GetUrls() {
		jsonData, _ := json.Marshal(w.sampleData)
		reader := bytes.NewReader(jsonData)
		go http.Post(url, "application/json", reader)
	}
}

func PeriodicalWarmUp() {
	for {
		time.Sleep(Delay)
		mutex.Lock()
		for _, warmer := range Manager {
			warmer.WarmUp()
		}
		mutex.Unlock()
	}
}
