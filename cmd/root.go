package cmd

import (
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	urlBase = "http://localhost:6000"
)

var rootCmd = &cobra.Command{
	Use: "drip",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

}
