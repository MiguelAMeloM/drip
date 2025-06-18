/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package modelProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MiguelAMeloM/drip/core/serverRaiser"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os/exec"
)

type ShadowProxy struct {
	ProxyBase
	portStable int
	portShadow int
}

func NewShadowProxy(ms NewModelSetting) Proxy {
	modelName := ms.ModelName
	raiser := ms.Raiser
	var portStable, portShadow int
	var err error
	portStable, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	portShadow, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	cmdStable, endpoint := raiser(portStable, modelName, serverRaiser.Prod)
	cmdShadow, _ := raiser(portShadow, modelName, serverRaiser.Dev)
	return &ShadowProxy{
		ProxyBase: ProxyBase{
			activeServers: []*exec.Cmd{cmdStable, cmdShadow},
			modelName:     modelName,
			endpoint:      endpoint,
		},
		portStable: portStable,
		portShadow: portShadow,
	}
}

func (p *ShadowProxy) ForwardRequest(request *http.Request) (gin.H, error) {
	var urlA, urlB string
	urlA = fmt.Sprintf("%s:%d/%s", baseUrl, p.portStable, p.endpoint)
	urlB = fmt.Sprintf("%s:%d/%s", baseUrl, p.portShadow, p.endpoint)

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	responseA, err := http.Post(urlA, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer responseA.Body.Close()

	dataA := make(gin.H)
	err = json.NewDecoder(responseA.Body).Decode(&dataA)
	if err != nil {
		return nil, err
	}

	responseB, err := http.Post(urlB, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer responseB.Body.Close()

	dataB := make(gin.H)
	err = json.NewDecoder(responseB.Body).Decode(&dataB)
	if err != nil {
		return nil, err
	}

	result := gin.H{
		"data": gin.H{
			"champion":   dataA,
			"challenger": dataB,
		},
		"metadata": gin.H{
			"reqId": genReqId(),
		},
	}

	return result, nil
}

func (p *ShadowProxy) Close() {
	shutDownServer(&p.ProxyBase)
}

func (p *ShadowProxy) GetUrls() []string {
	return []string{
		fmt.Sprintf("%s:%d/%s", baseUrl, p.portShadow, p.endpoint),
		fmt.Sprintf("%s:%d/%s", baseUrl, p.portStable, p.endpoint),
	}
}
