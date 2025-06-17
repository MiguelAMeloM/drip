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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	inHouseEndpoint = "invocations"
)

var inHouseFolder string = filepath.Join(os.Getenv("HOME"), "in_house_models")

func UpdateInHouseFolder(folder string) {
	inHouseFolder = filepath.Join(os.Getenv("HOME"), folder)
}

func InHouse(port int, modelName string, alias string) (*exec.Cmd, string) {
	path := inHouseFolder + "/" + modelName
	var arguments = []string{
		"python3",
		path,
		"--alias",
		alias,
		"--port",
		strconv.Itoa(port),
	}
	cmd := exec.Command("python", arguments...)

	outputChan := getOutputInterface(modelName)

	cmd.Stdout = outputChan
	cmd.Stderr = outputChan

	go raiseServer(cmd)

	return cmd, inHouseEndpoint
}
