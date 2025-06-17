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
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func removeModel(c *gin.Context) {
	r := c.Request
	if !isAllowed(r) {
		c.Error(fmt.Errorf("access denied"))
		return
	}
	modelName := r.URL.Query().Get("model")
	if modelName == "" {
		c.Error(fmt.Errorf("model name is required"))
		return
	}

	err := core.Manager.RemoveModel(modelName)
	if err != nil {
		c.Error(err)
		return
	} else {
		c.String(http.StatusOK, "Model removed")
	}
}
