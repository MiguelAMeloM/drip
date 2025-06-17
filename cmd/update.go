package cmd

import (
	"drip/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		cobra.CheckErr(err)

		modelName := viper.GetString("model")
		releaseType := viper.GetString("release_type")
		backendType := viper.GetString("backend_type")
		deadline := viper.GetString("deadline")
		prob := viper.GetString("prob")

		u, err := url.Parse(urlBase + "/" + engine.Upd)
		if err != nil {
			panic(err)
		}

		q := u.Query()
		q.Set("model", modelName)
		q.Set("release_type", releaseType)
		q.Set("backend_type", backendType)
		q.Set("deadline", deadline)
		q.Set("prob", prob)

		u.RawQuery = q.Encode()

		done := make(chan bool)
		go spinner(done, "updating model...")

		resp, err := http.Get(u.String())
		done <- true
		if err != nil {
			cobra.CheckErr(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Println(resp.Status)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("model", "m", "", "model name")
	updateCmd.Flags().StringP("release_type", "r", "stable", "release type")
	updateCmd.Flags().StringP("backend_type", "b", "mlflow", "backend type")
	updateCmd.Flags().StringP("deadline", "t", "", "deadline")
	updateCmd.Flags().StringP("prob", "p", "0.5", "prob")
}
