package serverRaiser

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	//stages
	Production  = "production"
	Development = "development"

	dir = "logs"
)

var (
	Prod = "prod"
	Dev  = "dev"
)

func UpdateAlias(stage string, alias string) {
	switch stage {
	case Production:
		Prod = alias
	case Development:
		Dev = alias
	default:
		panic("invalid stage value")
	}
}

type RaiseServer func(int, string, string) (*exec.Cmd, string)

func getOutputInterface(modelName string) io.Writer {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(fmt.Errorf("error creando directorio: %v", err))
		}
	}

	// Obtener fecha y hora actual formateada
	timestamp := time.Now().Format("20060102_150405")

	// Construir el nombre del archivo con timestamp
	fileName := fmt.Sprintf("%s_%s.log", modelName, timestamp)

	filePath := filepath.Join(dir, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("error abriendo archivo: %v", err))
	}

	return file
}

func raiseServer(cmd *exec.Cmd) {
	if err := cmd.Run(); err != nil {
		log.Println("raiseServer:", err)
	}
}
