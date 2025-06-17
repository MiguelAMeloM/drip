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
