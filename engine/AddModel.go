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
	"drip/core/modelProxy"
	"drip/core/serverRaiser"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func addModel(c *gin.Context) {
	r := c.Request

	if !isAllowed(r) {
		c.Error(fmt.Errorf("access denied"))
		return
	}
	modelName := r.URL.Query().Get("model")
	if modelName == "" {
		c.Error(fmt.Errorf("model name is required"))
		return
	} else {
		if _, ok := core.Manager.ModelsMapping[modelName]; ok {
			http.Error(c.Writer, fmt.Sprintf("model %s already exists", modelName), http.StatusBadRequest)
			return
		}
	}
	releaseType := r.URL.Query().Get("release_type")
	if releaseType == "" {
		c.Error(fmt.Errorf("release type is required"))
		return
	}
	backendType := r.URL.Query().Get("backend_type")
	if backendType == "" {
		c.Error(fmt.Errorf("backend type is required"))
		return
	}
	deadlineString := r.URL.Query().Get("deadline")
	var deadline time.Time
	if deadlineString != "" {
		var err error
		deadline, err = time.Parse("2006-01-02", deadlineString)
		if err != nil {
			c.Error(fmt.Errorf("deadline is invalid"))
			return
		}
	}
	probS := r.URL.Query().Get("prob")
	var prob float64
	if probS == "" && releaseType == core.Canary {
		c.Error(fmt.Errorf("prob is required"))
		return
	} else {
		var err error
		prob, err = strconv.ParseFloat(probS, 64)
		if err != nil {
			c.Error(fmt.Errorf("prob is invalid"))
			return
		}
	}

	modelSetting := modelProxy.NewModelSetting{
		ModelName:   modelName,
		ReleaseType: releaseType,
		Deadline:    deadline,
		Prob:        prob,
	}

	switch backendType {
	case "mlflow":
		modelSetting.Raiser = serverRaiser.RaiseMLFlowServer

	case "inhouse":
		modelSetting.Raiser = serverRaiser.InHouse
	default:
		c.Error(fmt.Errorf("backend type is invalid"))
		return
	}
	err := core.Manager.AddModel(modelSetting)
	if err != nil {
		c.Error(err)
		return
	}
	c.String(http.StatusCreated, "Model added")
}
