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
	"encoding/json"
	"fmt"
	"github.com/MiguelAMeloM/drip/core/serverRaiser"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os/exec"
)

type ABTesting struct {
	ProxyBase
	portA int
	portB int
	prob  float64
}

func NewABTestingProxy(ms NewModelSetting) Proxy {
	modelName := ms.ModelName
	raiser := ms.Raiser
	prob := ms.Prob
	var portA, portB int
	var err error
	portA, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	portB, err = getAvailablePort()
	if err != nil {
		panic(err)
	}
	cmdA, endpoint := raiser(portA, modelName, serverRaiser.Prod)
	cmdB, _ := raiser(portB, modelName, serverRaiser.Dev)

	return &ABTesting{
		ProxyBase: ProxyBase{
			activeServers: []*exec.Cmd{cmdA, cmdB},
			modelName:     modelName,
			endpoint:      endpoint,
		},
		portA: portA,
		portB: portB,
		prob:  prob,
	}
}

func (p *ABTesting) forwardReq(request *http.Request, useA bool) (gin.H, error) {
	var url, testCase string
	if useA {
		url = fmt.Sprintf("%s:%d/%s", baseUrl, p.portA, p.endpoint)
		testCase = "champion"
	} else {
		url = fmt.Sprintf("%s:%d/%s", baseUrl, p.portB, p.endpoint)
		testCase = "challenger"
	}

	response, err := http.Post(url, "application/json", request.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := make(gin.H)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	result := gin.H{
		"data": data,
		"metadata": map[string]string{
			"case":  testCase,
			"reqId": genReqId(),
		},
	}

	return result, nil
}

func (p *ABTesting) ForwardRequest(request *http.Request) (gin.H, error) {
	useA := rand.Float64() >= p.prob
	return p.forwardReq(request, useA)
}

func (p *ABTesting) Close() {
	shutDownServer(&p.ProxyBase)
}

func (p *ABTesting) GetUrls() []string {
	return []string{
		fmt.Sprintf("%s:%d/%s", baseUrl, p.portA, p.endpoint),
		fmt.Sprintf("%s:%d/%s", baseUrl, p.portB, p.endpoint),
	}
}
