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
	"github.com/gin-gonic/gin"
)

func UpdateModel(c *gin.Context) {
	removeModel(c)
	addModel(c)
}
