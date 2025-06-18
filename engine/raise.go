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
	"fmt"
	"github.com/MiguelAMeloM/drip/core"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zsais/go-gin-prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	Inf     = "inference"
	Add     = "add_model"
	Rem     = "remove_model"
	Upd     = "update_model"
	AliasCh = "change_alias"

	Ps = "list_models"
	SD = "shut_down"

	Stats      = "stats"
	Autoscaler = "autoscaler"
)

func inference(c *gin.Context) {
	r := c.Request
	w := c.Writer
	if !isAllowed(r) {
		c.Error(fmt.Errorf("not allowed"))
		return
	}

	modelName := r.URL.Query().Get("modelName")
	if modelName == "" {
		http.Error(w, "modelName is required", http.StatusBadRequest)
		return
	}

	resp, err := core.Manager.ForwardRequest(modelName, r)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, resp)
}

func shutDown(c chan os.Signal) {
	<-c
	for name, proxy := range core.Manager.ModelsMapping {
		fmt.Printf("shutting down model proxy %s\n", name)
		proxy.Close()
	}
	os.Exit(0)
}

func RaiseMainServer(token string) {
	AuthToken = token

	router := gin.Default()

	router.POST("/"+Inf, inference)
	router.GET("/"+Add, addModel)
	router.GET("/"+Rem, removeModel)
	router.GET("/"+AliasCh, UpdateStageAlias)
	router.GET("/"+Upd, UpdateModel)
	router.GET("/"+Ps, listModels)
	router.GET("/"+SD, shutDownEndpoint)
	router.GET("/"+Stats, viewStats)
	router.GET("/"+Autoscaler, UpdateAutoScaler)

	pr := ginprometheus.NewPrometheus("gin")
	pr.Use(router)

	p := viper.GetInt("port")
	port := fmt.Sprintf(":%d", p)
	fmt.Println("Serving on http://0.0.0.0" + port)
	fmt.Println("Press Ctrl-C to stop")
	fmt.Println("Authentication token: " + AuthToken)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go shutDown(sigChan)
	core.RaiseSubprocesses()

	err := router.Run(port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
