package engine

import (
	"github.com/gin-gonic/gin"
)

func UpdateModel(c *gin.Context) {
	removeModel(c)
	addModel(c)
}
