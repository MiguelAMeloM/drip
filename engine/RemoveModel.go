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
