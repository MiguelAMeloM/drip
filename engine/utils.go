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
	"drip/core"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

func listModels(c *gin.Context) {
	r := c.Request
	w := c.Writer
	if !isAllowed(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	models := make([]string, 0)
	for model, _ := range core.Manager.ModelsMapping {
		models = append(models, model)
	}
	sort.Strings(models)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(models)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(resp)
}

func shutDownEndpoint(c *gin.Context) {
	r := c.Request
	w := c.Writer
	if !isAllowed(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	go shutDownDelay()
}

func shutDownDelay() {
	time.Sleep(5 * time.Second)
	for name, proxy := range core.Manager.ModelsMapping {
		fmt.Printf("shutting down model proxy %s\n", name)
		proxy.Close()
	}
	os.Exit(0)
}

func viewStats(c *gin.Context) {
	r := c.Request
	w := c.Writer
	if !isAllowed(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	stats, err := core.Manager.GenStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateAutoScaler(c *gin.Context) {
	r := c.Request
	w := c.Writer
	if !isAllowed(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	secondsS := r.URL.Query().Get("seconds")
	if secondsS == "" {
		secondsS = "0"
	}
	secondsFloat, err := strconv.Atoi(secondsS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	seconds := time.Duration(secondsFloat) * time.Second
	minReq := r.URL.Query().Get("min")
	if minReq == "" {
		http.Error(w, "min is required", http.StatusBadRequest)
	}
	minReqFloat, err := strconv.ParseFloat(minReq, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	maxReq := r.URL.Query().Get("max")
	if maxReq == "" {
		http.Error(w, "max is required", http.StatusBadRequest)
	}
	maxReqFloat, err := strconv.ParseFloat(maxReq, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = core.UpdateAutoscalingParams(seconds, minReqFloat, maxReqFloat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
