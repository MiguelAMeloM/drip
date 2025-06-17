/*
 * Copyright (c) 2025.
 * Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 */

package serverRaiser

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strconv"
)

const (
	mlflowEndpoint = "invocations"
)

func RaiseMLFlowServer(port int, modelName string, alias string) (*exec.Cmd, string) {
	var arguments = []string{
		"models",
		"serve",
		"-m",
		fmt.Sprintf("models:/%s@%s", modelName, alias),
		"--port",
		strconv.Itoa(port),
		"--host",
		"0.0.0.0",
		"--env-manager",
		"local",
	}
	trackingUri := viper.GetString("mlflow.tracking_uri")
	cmd := exec.Command("mlflow", arguments...)
	cmd.Env = append(os.Environ(), "MLFLOW_TRACKING_URI="+trackingUri)

	outputChan := getOutputInterface(modelName)

	cmd.Stdout = outputChan
	cmd.Stderr = outputChan

	go raiseServer(cmd)

	return cmd, mlflowEndpoint
}
