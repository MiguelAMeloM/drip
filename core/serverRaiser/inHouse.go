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
