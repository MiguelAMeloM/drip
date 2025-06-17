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
	"drip/core/serverRaiser"
	"fmt"
	"github.com/gin-gonic/gin"
)

func UpdateStageAlias(c *gin.Context) {
	r := c.Request
	if !isAllowed(r) {
		c.Error(fmt.Errorf("access denied"))
		return
	}
	stage := r.URL.Query().Get("stage")
	if stage == "" {
		c.Error(fmt.Errorf("stage is required"))
	}
	alias := r.URL.Query().Get("alias")
	if alias == "" {
		c.Error(fmt.Errorf("alias is required"))
	}
	serverRaiser.UpdateAlias(stage, alias)
}
