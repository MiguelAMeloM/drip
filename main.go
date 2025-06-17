package main

import (
	"drip/cmd"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cmd.Execute()
}
