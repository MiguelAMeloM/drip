package cmd

import (
	"crypto/rand"
	"drip/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		cobra.CheckErr(err)
		token, err := cmd.Flags().GetString("token")
		cobra.CheckErr(err)
		engine.RaiseMainServer(token)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 6000, "port")
	serveCmd.Flags().String("token", rand.Text(), "Authentication token")
}
