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
