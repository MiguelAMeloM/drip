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
	"drip/core/serverRaiser"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

type StableRelease struct {
	ProxyBase
	port int
}

func NewStableRelease(ms NewModelSetting) Proxy {
	modelName := ms.ModelName
	raiser := ms.Raiser
	port, err := getAvailablePort()
	if err != nil {
		panic(err)
	}
	cmd, endpoint := raiser(port, modelName, serverRaiser.Prod)
	return &StableRelease{
		ProxyBase: ProxyBase{
			activeServers: []*exec.Cmd{cmd},
			modelName:     modelName,
			endpoint:      endpoint,
		},
		port: port,
	}
}

func (s *StableRelease) ForwardRequest(request *http.Request) (gin.H, error) {
	url := fmt.Sprintf("%s:%d/%s", baseUrl, s.port, s.endpoint)
	response, err := http.Post(url, "application/json", request.Body)
	if err != nil {
		return nil, err
	} else {
		data := make(gin.H)
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		result := gin.H{
			"data": data,
			"metadata": gin.H{
				"reqId": genReqId(),
			},
		}
		return result, nil
	}
}

func (s *StableRelease) Close() {
	shutDownServer(&s.ProxyBase)
}
